// Package repository provides MongoDB-backed persistence for message audit logs.
// Each message delivery attempt (success or failure) is stored in the "messages"
// collection for audit trail and debugging purposes.
package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"ub-communicator/pkg/messaging"
	"ub-communicator/pkg/platform"
)

const CollectionName = "messages"

type messageRepository struct {
	db     *mongo.Client
	dbName string
}

// NewMessage persists a message delivery record to MongoDB.
// Uses a 5-second timeout to prevent indefinite blocking.
func (mr *messageRepository) NewMessage(message *messaging.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mr.db.Database(mr.dbName).Collection(CollectionName)
	_, err := collection.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to insert message into %s.%s: %w", mr.dbName, CollectionName, err)
	}
	return nil
}

// NewMessageRepository creates a Repository backed by the given MongoDB client.
// The database name is read from the "mongodb.name" config key.
func NewMessageRepository(db *mongo.Client, c platform.Configs) messaging.Repository {
	dbName := c.GetString("mongodb.name")
	return &messageRepository{db, dbName}
}
