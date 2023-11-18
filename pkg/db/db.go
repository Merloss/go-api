package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
