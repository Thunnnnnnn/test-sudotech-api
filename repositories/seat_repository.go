package repositories

import (
	"context"
	"log"

	"gin-api/database"
	"gin-api/models"

	"go.mongodb.org/mongo-driver/bson"
)

func FindAllSeats() ([]models.Seat, error) {
	ctx := context.Background()

	if database.SeatCollection == nil {
		log.Fatal("SeatCollection is not initialized")
		return nil, nil
	}

	cursor, err := database.SeatCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var seats []models.Seat
	if err := cursor.All(ctx, &seats); err != nil {
		return nil, err
	}

	return seats, nil
}
