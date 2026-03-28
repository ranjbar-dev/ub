package messaging

import (
	"ub-communicator/pkg/platform"
)

// MailService wraps platform.MailerClient to provide a messaging-layer
// abstraction for email delivery. This wrapper exists as an extension point
// for adding cross-cutting concerns (logging, metrics, retry, circuit breaker)
// without modifying the platform-level mail provider implementations.
//
// Currently it delegates directly to the underlying MailerClient, but future
// enhancements should be added here rather than in the platform layer.
type MailService interface {
	// Send delivers an email with the given subject and HTML content to the receiver.
	// Returns (true, nil) on success, (false, error) on failure.
	Send(subject string, receiver string, content string) (bool, error)
}

// WHY: This wrapper delegates directly to MailerClient with no additional logic.
// It exists as an extension point for future cross-cutting concerns (logging,
// metrics, retry, circuit breaker) at the messaging layer without modifying
// the platform-level provider implementations.
type mailService struct {
	mc platform.MailerClient
}

func (s *mailService) Send(subject string, receiver string, content string) (bool, error) {
	return s.mc.Send(subject, receiver, content)
}

// NewMailService creates a MailService backed by the given MailerClient.
func NewMailService(mc platform.MailerClient) MailService {
	return &mailService{mc}
}
