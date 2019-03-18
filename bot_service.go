package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

type MailService struct {
	from       FromEmailAddress
	token      string
	httpClient *http.Client
	url        string
}

func NewMailService(httpClient *http.Client, token, from, url string) *MailService {
	return &MailService{
		from:       FromEmailAddress{Email: from},
		token:      token,
		httpClient: httpClient,
		url:        url,
	}
}

func getEmailFromText(text string) string {
	re := regexp.MustCompile("[a-z0-9-][a-z0-9-.]{1,30}@[a-z0-9-]{1,65}.[a-z]{1,}")
	return re.FindString(text)
}

func renderTemplate(filename string) (string, error) {
	var writer bytes.Buffer
	template := template.Must(template.ParseFiles(filename))
	err := template.Execute(&writer, nil)
	if err != nil {
		return "", errors.New("Can't execute the template: " + err.Error())
	}
	return writer.String(), nil
}

func (ms *MailService) MakeConfig(toEmail, subject, message, messageType string) MessageConfig {
	to := []ToEmailAddress{{Email: toEmail}}
	personalizations := []EmailPersonalization{{To: to, Subject: subject}}
	content := []Content{{Type: messageType, Value: message}}
	return MessageConfig{Personalizations: personalizations, From: ms.from, Content: content}
}

func (ms *MailService) SendEmail(msg *MessageConfig) error {

	jsonPayload := new(bytes.Buffer)
	err := json.NewEncoder(jsonPayload).Encode(msg)
	if err != nil {
		log.Print("[E]: When marshaling errorMessage, e:", err.Error())
		return err
	}

	req, err := http.NewRequest("POST", ms.url, jsonPayload)
	if err != nil {
		log.Print("Failed to make request to Sendgrid", err.Error())
		return err
	}
	req.Header.Add("Authorization", ms.token)
	req.Header.Add("content-type", "application/json")

	resp, err := ms.httpClient.Do(req)
	if err != nil {
		log.Print("Failed to make request to Sendgrid", err.Error())
		return err
	}

	defer resp.Body.Close()

	req.Body.Close()
	req.Close = true

	return nil
}
