package messaging

import "time"

const (
	MessageStatusPending    = "pending"
	MessageStatusFailed     = "failed"
	MessageStatusSuccessful = "successful"
)

const (
	MessageTypeEmail = "EMAIL"
	MessageTypeSms   = "SMS"
)

// Message represents a notification to be delivered via email or SMS.
// It is unmarshaled from the RabbitMQ message body (JSON) and persisted
// to MongoDB after delivery attempt (with updated Status).
type Message struct {
	Receiver    string    `bson:"receiver" json:"receiver"`
	Subject     string    `bson:"subject" json:"subject"`
	Content     string    `bson:"content" json:"content"`
	// TODO: Priority is persisted to MongoDB but never read back or used to influence delivery order.
	Priority int `bson:"priority" json:"priority"`
	// TODO: ScheduledAt is persisted to MongoDB but never read back or used to delay/schedule delivery.
	ScheduledAt string `bson:"scheduledAt" json:"scheduledAt"`
	Type        string    `bson:"type" json:"type"`
	Status      string    `bson:"status" json:"status"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
}

// Repository persists message delivery audit logs to MongoDB.
type Repository interface {
	// NewMessage stores a completed message record (sent or failed).
	// Called after every delivery attempt regardless of outcome.
	NewMessage(message *Message) error
}
