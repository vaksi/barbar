package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connect is a func to connect client mongodbs
func Connect(ctx context.Context, url, database string) (*mongo.Database, error) {

	clientOptions := options.Client()
	clientOptions.ApplyURI(url)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// check ping connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	conn := client.Database(database)

	return conn, nil
}
