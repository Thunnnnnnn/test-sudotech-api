package repositories

import (
	"context"
	"log"

	"gin-api/database"
	"gin-api/models"

	"go.mongodb.org/mongo-driver/bson"
)

func FindAllUsers() ([]models.User, error) {
	ctx := context.Background()

	if database.UserCollection == nil {
		log.Fatal("UserCollection is not initialized")
		return nil, nil
	}

	cursor, err := database.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
