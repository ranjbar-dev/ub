package platform

import (
	"github.com/mailjet/mailjet-apiv3-go/v3"
	"strings"
)

type mailJetMailerClient struct {
	mailerClient *mailjet.Client
	name         string
	fromAddress  string
	logger       Logger
}

func (m *mailJetMailerClient) Send(subject string, receiver string, content string) (bool, error) {

	plainTextContent := subject
	// WHY: Subject is prefixed with [UNITEDBIT] to identify automated emails
	// in the recipient's inbox. The check for "[" prevents double-prefixing
	// if the subject already contains a tag.
	if !strings.Contains(subject, "[") {
		subject = "[" + m.name + "] " + subject
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: m.fromAddress,
				Name:  m.name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: receiver,
					//Name:  "passenger 1",
				},
			},
			Subject:  subject,
			TextPart: plainTextContent,
			HTMLPart: content,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := m.mailerClient.SendMailV31(&messages)
	if err != nil {
		return false, err
	}

	return true, nil

}
