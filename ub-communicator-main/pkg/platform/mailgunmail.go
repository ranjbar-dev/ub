package platform

import (
	"github.com/mailgun/mailgun-go"
)

type mailGunClient struct {
	mailerClient *mailgun.MailgunImpl
	name         string
	fromAddress  string
	logger       Logger
}

func (m *mailGunClient) Send(subject string, receiver string, content string) (bool, error) {

	plainTextContent := subject
	message := m.mailerClient.NewMessage(m.fromAddress, subject, plainTextContent, receiver)
	message.SetHtml(content)

	_, _, err := m.mailerClient.Send(message)
	if err != nil {
		return false, err
	}

	return true, nil
}
