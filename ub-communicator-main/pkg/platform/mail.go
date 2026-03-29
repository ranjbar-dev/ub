package platform

import (
	"crypto/tls"
	"github.com/mailgun/mailgun-go"
	"github.com/mailjet/mailjet-apiv3-go/v3"
	"github.com/sendgrid/sendgrid-go"
	"gopkg.in/gomail.v2"
)

const (
	// MailerSendGrid selects SendGrid as the email delivery provider.
	MailerSendGrid = "sendgrid"
	// MailerMailJet selects Mailjet as the email delivery provider.
	MailerMailJet = "mailjet"
	// MailerMailGun selects Mailgun as the email delivery provider.
	MailerMailGun = "mailgun"
	// MailerSMTP selects direct SMTP as the email delivery provider.
	MailerSMTP = "smtp"
)

// MailerClient is the interface implemented by all email delivery providers.
// Implementations: sendGridMailerClient, mailJetMailerClient, mailGunClient, smtpClient.
type MailerClient interface {
	// Send delivers an email with the given subject and HTML content to the receiver.
	// Returns (true, nil) on success, (false, error) on failure.
	Send(subject string, receiver string, content string) (bool, error)
}

// NewMailerClient is a factory that creates the appropriate MailerClient
// based on the "mailer_broker" config value ("sendgrid", "mailjet", "mailgun", "smtp").
// Returns nil if the value is unrecognized.
func NewMailerClient(configs Configs, logger Logger) MailerClient {
	name := configs.GetString("mail.name")
	fromAddress := configs.GetString("mail.from_address")

	mailerBroker := configs.GetString("mailer_broker")

	switch mailerBroker {
	case MailerSendGrid:
		apiKey := configs.GetString("sendgrid.api_key")
		client := sendgrid.NewSendClient(apiKey)

		return &sendGridMailerClient{client, name, fromAddress, logger}

	case MailerMailJet:
		apiPublicKey := configs.GetString("mailjet.api_public_key")
		apiPrivateKey := configs.GetString("mailjet.api_private_key")

		client := mailjet.NewMailjetClient(apiPublicKey, apiPrivateKey)
		return &mailJetMailerClient{client, name, fromAddress, logger}

	case MailerMailGun:
		apiKey := configs.GetString("mailgun.api_key")
		domain := configs.GetString("mailgun.domain")
		apiBase := configs.GetString("mailgun.api_base")

		client := mailgun.NewMailgun(domain, apiKey)
		client.SetAPIBase(apiBase)

		return &mailGunClient{client, name, fromAddress, logger}

	case MailerSMTP:
		client := gomail.NewDialer(configs.GetString("smtp.host"), configs.GetInt("smtp.port"),
			configs.GetString("smtp.username"), configs.GetString("smtp.password"))

		client.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		}

		return &smtpClient{client, name, fromAddress, logger}

	default:
		logger.Error("unknown mailer_broker value: " + mailerBroker)
		return nil
	}

}
