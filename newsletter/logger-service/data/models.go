package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// This function is used to create a new mongo client
func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

// This function is used to insert the log recieved in database
func Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}
