// Package messaging orchestrates message delivery by routing messages to the
// appropriate channel (email or SMS) based on the message Type field.
// It handles message creation from raw JSON, validation, delivery via
// the configured provider, and persistence of audit logs to MongoDB.
//
// Key types:
//   - Service: routes messages and persists delivery results
//   - MailService: email delivery abstraction (wraps platform.MailerClient)
//   - SmsService: SMS delivery via Twilio HTTP API
//   - Repository: MongoDB persistence for message audit logs
//   - Message: the data model (JSON from RabbitMQ ↔ BSON in MongoDB)
package messaging

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"ub-communicator/pkg/platform"

	"go.uber.org/zap"
)

// Service orchestrates message delivery by routing to email or SMS
// and persisting audit logs to MongoDB.
type Service interface {
	// Send routes the message to the appropriate delivery channel based on Type,
	// updates the message Status, and persists the result to MongoDB.
	// Precondition: message.Type must be "EMAIL" or "SMS" (uppercase).
	Send(message Message) error
	// CreateMessage parses raw JSON bytes from RabbitMQ into a Message struct.
	// It normalizes the Type field to uppercase and sets initial Status and CreatedAt.
	CreateMessage(data []byte) (Message, error)
}

type service struct {
	mr     Repository
	ms     MailService
	ss     SmsService
	logger platform.Logger
}

var phoneRegexp = regexp.MustCompile(`^\+[1-9]\d{6,14}$`)

// validateMessage checks that a message has the required fields and valid format.
// Returns an error describing the first validation failure found.
func validateMessage(msg Message) error {
	if msg.Receiver == "" {
		return fmt.Errorf("receiver is empty")
	}
	if msg.Content == "" {
		return fmt.Errorf("content is empty")
	}

	switch msg.Type {
	case MessageTypeEmail:
		if _, err := mail.ParseAddress(msg.Receiver); err != nil {
			return fmt.Errorf("invalid email address %q: %w", msg.Receiver, err)
		}
	case MessageTypeSms:
		if !phoneRegexp.MatchString(msg.Receiver) {
			return fmt.Errorf("invalid phone number %q: must be E.164 format (+<country><number>)", msg.Receiver)
		}
	default:
		return fmt.Errorf("unknown message type %q: expected %q or %q", msg.Type, MessageTypeEmail, MessageTypeSms)
	}
	return nil
}

func (s *service) Send(message Message) error {
	if err := validateMessage(message); err != nil {
		message.Status = MessageStatusFailed
		if saveErr := s.mr.NewMessage(&message); saveErr != nil {
			s.logger.Error("failed to save invalid message to db", zap.Error(saveErr))
		}
		return fmt.Errorf("message validation failed: %w", err)
	}

	var deliveryErr error

	switch message.Type {
	case MessageTypeEmail:
		isSent, err := s.ms.Send(message.Subject, message.Receiver, message.Content)
		if isSent && err == nil {
			message.Status = MessageStatusSuccessful
		} else {
			if err != nil {
				s.logger.Error("failed to send email", zap.Error(err))
				deliveryErr = fmt.Errorf("email delivery failed: %w", err)
			} else {
				deliveryErr = fmt.Errorf("email delivery failed: provider returned false with no error")
				s.logger.Error("email delivery failed: provider returned false with no error")
			}
			message.Status = MessageStatusFailed
		}
		if saveErr := s.mr.NewMessage(&message); saveErr != nil {
			s.logger.Error("failed to save message to database", zap.Error(saveErr))
		}

	case MessageTypeSms:
		isSent, err := s.ss.Send(message.Subject, message.Receiver, message.Content)
		if isSent && err == nil {
			message.Status = MessageStatusSuccessful
		} else {
			if err != nil {
				s.logger.Error("failed to send sms", zap.Error(err))
				deliveryErr = fmt.Errorf("sms delivery failed: %w", err)
			} else {
				deliveryErr = fmt.Errorf("sms delivery failed: provider returned false with no error")
				s.logger.Error("sms delivery failed: provider returned false with no error")
			}
			message.Status = MessageStatusFailed
		}
		if saveErr := s.mr.NewMessage(&message); saveErr != nil {
			s.logger.Error("failed to save message to database", zap.Error(saveErr))
		}
	}
	return deliveryErr
}

func (s *service) CreateMessage(data []byte) (Message, error) {
	var message Message
	if err := json.Unmarshal(data, &message); err != nil {
		return Message{}, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	message.Type = strings.ToUpper(message.Type)
	message.Status = MessageStatusPending
	message.CreatedAt = time.Now()

	return message, nil
}

// NewMessagingService creates a Service wired to the given repository, mail, SMS, and logger.
func NewMessagingService(mr Repository, ms MailService, ss SmsService, logger platform.Logger) Service {
	return &service{mr, ms, ss, logger}
}
