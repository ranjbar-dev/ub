// Package repository provides MongoDB-backed persistence for message audit logs.
// Each message delivery attempt (success or failure) is stored in the "messages"
// collection for audit trail and debugging purposes.
package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

// EnsureIndexes creates required indexes on the messages collection.
// Safe to call on every startup — CreateIndexes is idempotent for existing indexes.
func EnsureIndexes(ctx context.Context, db *mongo.Client, dbName string) error {
	coll := db.Database(dbName).Collection(CollectionName)
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "createdAt", Value: 1}}},
		{Keys: bson.D{{Key: "receiver", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
	}
	_, err := coll.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to ensure indexes on %s.%s: %w", dbName, CollectionName, err)
	}
	return nil
}
