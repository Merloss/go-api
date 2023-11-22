package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect establishes a connection to a MongoDB database and returns a pointer to the specified database.
// It uses the provided URI to connect to the MongoDB server and the specified database name.
// The function creates a context with a timeout of 15 seconds for the connection attempt.
//
// Usage:
//
//	db, err := Connect("mongodb://localhost:27017", "mydatabase")
func Connect(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return db, nil
}
