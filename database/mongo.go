package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection
var SeatCollection *mongo.Collection

func ConnectMongo() error {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is empty")
	}

	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(mongoURI),
	)
	if err != nil {
		return err
	}

	UserCollection = client.Database("testdb").Collection("users")
	SeatCollection = client.Database("testdb").Collection("seats")

	return nil
}
