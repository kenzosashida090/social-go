package mailer

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/kenzosashida090/social/internal/env"
)

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}
type ValidationEmailTemplate struct {
	UserName            string
	ActivateUrl         string
	GetStartedLink      string
	OnboardingVideoLink string
}
type MailtrapBody struct {
	fromEmail string
	to        string
	addrs     string
	client    smtp.Auth
}

func ConnectMailTrap() *MailtrapBody {
	username := env.GetString("MAILTRAP_USERNAME", "")
	password := env.GetString("MAILTRAP_PASSWORD", "")
	host := env.GetString("MAILTRAP_HOST", "")

	auth := smtp.PlainAuth("", username, password, host)

	return &MailtrapBody{
		fromEmail: "hello@example.com",
		addrs:     host + ":2525",
		client:    auth,
	}

}
func CreateTemplate(templ string, data any, nameTemplate string, emailbody *MailtrapBody, subject string) ([]byte, error) {
	t, err := template.New(nameTemplate).Parse(templ)
	if err != nil {
		fmt.Println(err)
	}

	var body bytes.Buffer
	err = t.Execute(&body, data)
	if err != nil {
		return []byte(""), err
	}
	return []byte(
		"FROM: " + emailbody.fromEmail + "\r\n" +
			emailbody.to +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body.String(),
	), nil

}
func (m *MailtrapBody) sendMail(to string, msg []byte) error {
	return smtp.SendMail(
		m.addrs,
		m.client,
		m.fromEmail,
		[]string{to},
		msg,
	)
}
func (m *MailtrapBody) SendActivationMail(token string, to string, username string, urlActivation string) error {
	templ := EmailValidationTemplate
	fmt.Println(urlActivation)
	data := &ValidationEmailTemplate{
		UserName:    username,
		ActivateUrl: urlActivation,
	}
	m.to = "To: " + to + "\r\n"
	msg, err := CreateTemplate(templ, data, "email", m, "Validate your account")
	if err != nil {
		return err
	}
	return m.sendMail(to, msg)
}
