package main

type SlackRequest struct {
	Token       string   `json:"token"`
	TeamID      string   `json:"team_id"`
	APIAppID    string   `json:"api_app_id"`
	Event       Event    `json:"event"`
	Type        string   `json:"type"`
	AuthedTeams []string `json:"authed_teams"`
	EventID     string   `json:"event_id"`
	EventTime   int64    `json:"event_time"`
	Challenge   string   `json:"challenge"`
}

type Event struct {
	Type        string `json:"type"`
	Channel     string `json:"channel"`
	User        string `json:"user"`
	Text        string `json:"text"`
	Ts          string `json:"ts"`
	EventTs     string `json:"event_ts"`
	ChannelType string `json:"channel_type"`
}

type MessageConfig struct {
	Personalizations []EmailPersonalization `json:"personalizations"`
	From             FromEmailAddress       `json:"from"`
	Content          []Content              `json:"content"`
}

type Content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type FromEmailAddress struct {
	Email string `json:"email"`
}

type EmailPersonalization struct {
	To      []ToEmailAddress `json:"to"`
	Subject string           `json:"subject"`
}

type ToEmailAddress struct {
	Email string `json:"email"`
}

type MessageDetails struct {
	email       string
	subgect     string
	message     string
	typeMessage string
}

type Config struct {
	Token                    string `yaml:"token"`
	URL                      string `yaml:"url"`
	Port                     int    `yaml:"port"`
	FromEmail                string `yaml:"from_email"`
	Subject                  string `yaml:"subject"`
	SubjectAboutWrongMessage string `yaml:"subject_about_wrong_message"`
	TypeMessage              string `yaml:"type_message"`
	Channel                  string `yaml:"channel"`
	EmailTemplate            string `yaml:"email_template"`
}
