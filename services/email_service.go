package services

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
	"os"
)

type EmailSender interface {
	SendEmail(to []*mail.Email, contentType, subject, body string) (*rest.Response, error)
}

type SendGridEmailSender struct {
	senderEmail string
	senderName  string
	apiKey      string
}

func NewSendGridEmailSender(config *viper.Viper) EmailSender {
	return &SendGridEmailSender{
		config.GetString("sendgrid.address"),
		config.GetString("sendgrid.name"),
		os.Getenv("SENDGRID_API_KEY"),
	}
}

func (s *SendGridEmailSender) SendEmail(to []*mail.Email, contentType, subject, body string) (*rest.Response, error) {
	from := mail.NewEmail(s.senderName, s.senderEmail)
	content := mail.NewContent(contentType, body)

	// Setup mail
	m := mail.NewV3Mail()
	m.SetFrom(from)
	m.Subject = subject
	p := mail.NewPersonalization()
	for _, recipient := range to {
		p.AddTos(recipient)
	}
	m.AddPersonalizations(p)
	m.AddContent(content)

	//Send it
	request := sendgrid.GetRequest(s.apiKey, "/v3/mail/send", "")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	return response, err
}
