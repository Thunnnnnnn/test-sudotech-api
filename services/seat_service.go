package services

import (
	"gin-api/models"

	"gin-api/repositories"

	"github.com/gin-gonic/gin"
)

func GetSeats() ([]models.Seat, error) {
	return repositories.FindAllSeats()
}

func CreateSeat(seat models.Seat) (models.Seat, error) {
	return repositories.CreateSeat(seat)
}

func BookSeat(c *gin.Context, seatID string) error {
	return repositories.BookSeat(c, seatID)
}

func CancelSeatBooking(c *gin.Context, seatID string) error {
	return repositories.CancelSeat(c, seatID)
}

func ConfirmSeatBooking(c *gin.Context, seatID string) error {
	return repositories.ConfirmSeat(c, seatID)
}
