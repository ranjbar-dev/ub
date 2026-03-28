package platform

import (
	"gopkg.in/gomail.v2"
)

type smtpClient struct {
	mailerClient *gomail.Dialer
	name        string
	fromAddress string
	logger      Logger
}

func (m *smtpClient) Send(subject string, receiver string, content string) (bool, error) {

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromAddress)
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", content)

	// Send the email
	if err := m.mailerClient.DialAndSend(message); err != nil {
		return false, err
	}

	return true, nil
}
