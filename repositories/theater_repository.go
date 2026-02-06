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

func FindAllTheaters() ([]models.Theater, error) {
	ctx := context.Background()

	if database.TheaterCollection == nil {
		log.Fatal("TheaterCollection is not initialized")
		return nil, nil
	}

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "seats"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "theater_id"},
			{Key: "as", Value: "seats"},
		}}},
	}

	cursor, err := database.TheaterCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var theaters []models.Theater
	if err := cursor.All(ctx, &theaters); err != nil {
		return nil, err
	}

	return theaters, nil
}

func CreateTheater(theater models.Theater) (models.Theater, error) {
	ctx := context.Background()

	if database.TheaterCollection == nil {
		log.Fatal("TheaterCollection is not initialized")
		return models.Theater{}, nil
	}

	result, err := database.TheaterCollection.InsertOne(ctx, theater)
	if err != nil {
		return models.Theater{}, err
	}

	theater.ID = result.InsertedID.(primitive.ObjectID)
	return theater, nil
}

func FindTheaterByID(id string) (models.Theater, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Theater{}, err
	}

	ctx := context.Background()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"_id": objID,
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "seats",      // collection seats
			"localField":   "_id",        // theater._id
			"foreignField": "theater_id", // seats.theater_id
			"as":           "seats",
		}}},
	}

	cursor, err := database.TheaterCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return models.Theater{}, err
	}
	defer cursor.Close(ctx)

	var theaters []models.Theater
	if err := cursor.All(ctx, &theaters); err != nil {
		return models.Theater{}, err
	}

	if len(theaters) == 0 {
		return models.Theater{}, mongo.ErrNoDocuments
	}

	return theaters[0], nil
}
