package mail

import (
	"net/smtp"
	"ticker-tracer/client/notification/mail/model"
)

type MailClientInterface interface {
	SendEmail(config model.EmailConfig, email model.Email) error
}
type MailHttpClient struct {
}

var mailHttpClient *MailHttpClient

func GetMailHttpClientInstance() *MailHttpClient {
	if mailHttpClient == nil {
		mailHttpClient = NewMailHttpClient()
	}
	return mailHttpClient
}

func NewMailHttpClient() *MailHttpClient {
	return &MailHttpClient{}
}

func (c *MailHttpClient) SendEmail(email model.Email) error {
	auth := smtp.PlainAuth("", "ticker.tracker.system@gmail.com", "nhpa axli mwii utbj", "smtp.gmail.com")
	addr := "smtp.gmail.com:587"
	subject := "Subject: " + email.Subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + mime + email.Body)

	return smtp.SendMail(addr, auth, "ticker.tracker.system@gmail.com", []string{email.To}, msg)
}
