package platform

// TODO: mailgun-go v2 (github.com/mailgun/mailgun-go v2.0.0+incompatible) is deprecated.
// Upgrade to v4: replace github.com/mailgun/mailgun-go with github.com/mailgun/mailgun-go/v4
// and update all import paths. The v4 API uses context-aware Send(ctx, message) instead of Send(message).
import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go"
)

const mailgunTimeout = 30 * time.Second

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

	ctx, cancel := context.WithTimeout(context.Background(), mailgunTimeout)
	defer cancel()

	type sendResult struct{ err error }
	done := make(chan sendResult, 1)
	go func() {
		_, _, err := m.mailerClient.Send(message)
		done <- sendResult{err: err}
	}()

	select {
	case <-ctx.Done():
		return false, fmt.Errorf("mailgun send timed out after %s: %w", mailgunTimeout, ctx.Err())
	case r := <-done:
		if r.err != nil {
			return false, r.err
		}
		return true, nil
	}
}
