package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type BotController struct {
	mailService *MailService
	mailConfig  *Config
}

func (bc *BotController) proccessSlack(w http.ResponseWriter, r *http.Request) {
	var req SlackRequest
	var msgConfig MessageConfig

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&req)

	switch req.Type {
	case "url_verification":
		io.WriteString(w, req.Challenge)
		break
	case "event_callback":
		if req.Event.Channel != bc.mailConfig.Channel {
			w.WriteHeader(http.StatusOK)
			break
		}
		email := getEmailFromText(req.Event.Text)
		message, err := renderTemplate(bc.mailConfig.EmailTemplate)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			log.Print("Can't execute the template" + err.Error())
			break
		}

		if email != "" {
			msgConfig = bc.mailService.MakeConfig(
				email,
				bc.mailConfig.Subject,
				message,
				bc.mailConfig.TypeMessage,
			)
		} else {
			msgConfig = bc.mailService.MakeConfig(
				bc.mailConfig.FromEmail,
				bc.mailConfig.SubjectAboutWrongMessage,
				req.Event.Text,
				bc.mailConfig.TypeMessage,
			)
		}
		bc.mailService.SendEmail(&msgConfig)
		w.WriteHeader(http.StatusOK)
		break
	}
}
