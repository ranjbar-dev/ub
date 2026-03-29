package platform

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"strings"
)

type sendGridMailerClient struct {
	mailerClient *sendgrid.Client
	name         string
	fromAddress  string
	logger       Logger
}

func (m *sendGridMailerClient) Send(subject string, receiver string, content string) (bool, error) {
	from := mail.NewEmail(m.name, m.fromAddress)
	to := mail.NewEmail("", receiver)
	plainTextContent := subject
	// WHY: Subject is prefixed with [UNITEDBIT] to identify automated emails
	// in the recipient's inbox. The check for "[" prevents double-prefixing
	// if the subject already contains a tag.
	if !strings.Contains(subject, "[") {
		subject = "[" + m.name + "] " + subject
	}
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, content)
	response, err := m.mailerClient.Send(message)

	if err != nil {
		return false, err
	} else {
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			return false, fmt.Errorf("unexpected status code from SendGrid: %d", response.StatusCode)
		}
		return true, nil
	}
}
