package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Ctx = context.Background()

func Connect() (*mongo.Database, error) {
	clientOptions := options.Client()
	mongoUri := fmt.Sprintf("mongodb://%s:%s", os.Getenv("LOCAL_MONGO_HOST"), os.Getenv("LOCAL_MONGO_PORT"))
	clientOptions.ApplyURI(mongoUri)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(Ctx)
	if err != nil {
		return nil, err
	}

	return client.Database("kawulo"), nil
}
