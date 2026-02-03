package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"gin-api/database"
	"gin-api/models"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func CreateSeat(seat models.Seat) (models.Seat, error) {
	_, err := database.SeatCollection.InsertOne(
		context.Background(),
		bson.M{
			"name":         seat.Name,
			"row":          seat.Row,
			"col":          seat.Col,
			"already_sold": seat.AlreadySold,
			"booked_at":    time.Now(),
			"expired_at":   time.Now(),
		},
	)

	return seat, err
}

func AcquireLock(key string, ttl time.Duration) (string, error) {
	rdb := database.RDB
	token := uuid.NewString()

	ok, err := rdb.SetNX(context.Background(), key, token, ttl).Result()

	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("lock already exists")
	}

	return token, nil
}

func FindByID(seatID string) (models.Seat, error) {
	objID, err := primitive.ObjectIDFromHex(seatID)
	if err != nil {
		return models.Seat{}, err
	}
	var seat models.Seat
	err = database.SeatCollection.FindOne(
		context.Background(),
		bson.M{"_id": objID},
	).Decode(&seat)
	return seat, err
}

var unlockScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`)

func ReleaseLock(key, token string) {
	rdb := database.RDB
	unlockScript.Run(context.Background(), rdb, []string{key}, token)
}

func ReleaseSeat(seatID string, lockKey string, token string) error {
	objectID, err := primitive.ObjectIDFromHex(seatID)
	rdb := database.RDB
	ttl, _ := rdb.TTL(context.Background(), lockKey).Result()

	if ttl <= 0 {
		_, err = database.SeatCollection.UpdateByID(
			context.Background(),
			objectID,
			bson.M{"$set": bson.M{"is_booked": false}},
		)

		return err
	} else {
		return errors.New("มีคนกำลังจองอยู่")
	}
}

func BookSeat(seatID string) error {
	lockKey := "lock:seat:" + seatID

	_, err := AcquireLock(lockKey, 5*time.Minute)

	if err != nil {
		return errors.New("มีคนกำลังจองอยู่")
	}

	seat, err := FindByID(seatID)
	if err != nil {
		return err
	}

	objID, _ := primitive.ObjectIDFromHex(seatID)

	if time.Now().After(*seat.ExpiredAt) {
		database.SeatCollection.UpdateByID(
			context.Background(),
			objID,
			bson.M{
				"$set": bson.M{
					"is_booked": false,
				},
				"$unset": bson.M{
					"booked_at":  "",
					"expired_at": "",
				},
			},
		)
	}
	now := time.Now()
	expire := now.Add(1 * time.Minute)
	_, err = database.SeatCollection.UpdateByID(
		context.Background(),
		objID,
		bson.M{"$set": bson.M{"booked_at": time.Now(), "expired_at": expire}},
	)

	return err
}
