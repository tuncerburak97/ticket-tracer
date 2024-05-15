package model

type EmailConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}
type Email struct {
	To      string
	Subject string
	Body    string
}
