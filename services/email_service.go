package services

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/dernise/base-api/models"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
)

type EmailSender interface {
	SendEmail(to []*mail.Email, contentType, subject, body string) (*rest.Response, error)
	SendEmailFromTemplate(user *models.User, subject string, templateLink string) (*rest.Response, error)
}

type SendGridEmailSender struct {
	senderEmail string
	senderName  string
	apiKey      string
	baseUrl     string
}

func NewSendGridEmailSender(config *viper.Viper) EmailSender {
	return &SendGridEmailSender{
		config.GetString("sendgrid_address"),
		config.GetString("sendgrid_name"),
		config.GetString("sendgrid_api_key"),
		config.GetString("base_url"),
	}
}

func (s SendGridEmailSender) SendEmail(to []*mail.Email, contentType, subject, body string) (*rest.Response, error) {
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
	return sendgrid.API(request)
}

type Data struct {
	User        *models.User
	HostAddress string
	AppName     string
}

func (s SendGridEmailSender) SendEmailFromTemplate(user *models.User, subject string, templateLink string) (*rest.Response, error) {

	to := mail.NewEmail(user.Firstname, user.Email)

	file, err := ioutil.ReadFile(templateLink)
	if err != nil {
		return nil, err
	}

	htmlTemplate := template.Must(template.New("emailTemplate").Parse(string(file)))

	data := Data{User: user, HostAddress: s.baseUrl, AppName: s.senderName}
	buffer := new(bytes.Buffer)
	err = htmlTemplate.Execute(buffer, data)
	if err != nil {
		return nil, err
	}

	return s.SendEmail([]*mail.Email{to}, "text/html", subject, buffer.String())
}
