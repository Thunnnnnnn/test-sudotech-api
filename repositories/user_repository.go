package repositories

import (
	"context"
	"log"

	"gin-api/database"
	"gin-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindAllUsers() ([]models.User, error) {
	ctx := context.Background()

	if database.UserCollection == nil {
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

func FindUserByEmail(email string) (models.User, error) {
	ctx := context.Background()

	if database.UserCollection == nil {
		log.Fatal("UserCollection is not initialized")
		return models.User{}, nil
	}
	var user models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return models.User{}, nil
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func CreateUser(user models.User) (models.User, error) {
	ctx := context.Background()

	if database.UserCollection == nil {
		log.Fatal("UserCollection is not initialized")
		return models.User{}, nil
	}

	result, err := database.UserCollection.InsertOne(ctx, user)
	if err != nil {
		panic(err)
		// return models.User{}, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}
