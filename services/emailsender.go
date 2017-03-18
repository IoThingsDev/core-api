package services

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sendgrid/rest"
	"github.com/spf13/viper"
)

type EmailSender interface {
	SendEmail(to []*mail.Email, contentType, subject, body string) (*rest.Response, error)
}

type SendGridEmailSender struct {
	senderEmail string
	senderName  string
	apiKey 	    string
}

func NewSendGridEmailSender(config *viper.Viper) EmailSender {
	return &SendGridEmailSender{
		config.GetString("sendgrid.address"),
		config.GetString("sendgrid.name"),
		config.GetString("sendgrid.apiKey"),
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
	for _,receipient := range to {
		p.AddTos(receipient)
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