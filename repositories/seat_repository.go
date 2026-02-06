package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gin-api/database"
	"gin-api/models"

	"github.com/gin-gonic/gin"
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
	fmt.Println("Creating seat:", seat)
	theaterID, _ := primitive.ObjectIDFromHex(seat.TheaterID.Hex())

	fmt.Println("Theater ID:", theaterID)
	_, err := database.SeatCollection.InsertOne(
		context.Background(),
		bson.M{
			"name":         seat.Name,
			"row":          seat.Row,
			"col":          seat.Col,
			"already_sold": seat.AlreadySold,
			"booked_at":    time.Now(),
			"expired_at":   time.Now(),
			"theater_id":   theaterID,
		},
	)

	return seat, err
}

func AcquireLock(key string, ttl time.Duration, userID string) (string, error) {
	ctx := context.Background()

	value := fmt.Sprintf("%s", userID)

	ok, err := database.RDB.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("seat is already locked")
	}

	return value, nil
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

func ReleaseSeat(seatID string, userID string) error {
	key := "lock:seat:" + seatID

	ctx := context.Background()
	val, err := database.RDB.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	expected := fmt.Sprintf("%s", userID)

	if val != expected {
		return errors.New("not lock owner")
	}

	return database.RDB.Del(ctx, key).Err()
}

func BookSeat(c *gin.Context, seatID string) error {
	userIDValue, exists := c.Get("user_id")

	if !exists {
		return errors.New("user not found in context")
	}

	haveSeat := ReleaseSeat(seatID, userIDValue.(string))

	if haveSeat == nil {
		return errors.New("คุณจองเก้าอี้นี้อยู่")
	}

	lockKey := "lock:seat:" + seatID
	_, err := AcquireLock(lockKey, 5*time.Minute, userIDValue.(string))

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
					"booked_at":  nil,
					"expired_at": nil,
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

func ConfirmSeat(c *gin.Context, seatID string) error {
	userIDValue, exists := c.Get("user_id")

	if !exists {
		return errors.New("user not found in context")
	}

	lockKey := "lock:seat:" + seatID + ":user_id:" + userIDValue.(string)

	objectID, err := primitive.ObjectIDFromHex(seatID)
	if err != nil {
		return err
	}

	_, err = database.SeatCollection.UpdateByID(
		context.Background(),
		objectID,
		bson.M{"$set": bson.M{"already_sold": true}},
	)

	ReleaseLock(lockKey, "")

	return err
}

func CancelSeat(c *gin.Context, seatID string) error {
	userIDValue, exists := c.Get("user_id")

	if !exists {
		return errors.New("user not found in context")
	}

	lockKey := "lock:seat:" + seatID + ":user_id:" + userIDValue.(string)

	ReleaseLock(lockKey, "")

	objID, _ := primitive.ObjectIDFromHex(seatID)

	// seat, err := FindByID(seatID)

	// now := time.Now()
	// expire := now.Add(1 * time.Minute)
	_, err := database.SeatCollection.UpdateByID(
		context.Background(),
		objID,
		bson.M{"$set": bson.M{"booked_at": time.Now().Add(time.Minute * -5), "expired_at": time.Now().Add(time.Minute * -5)}},
	)

	return err
}
