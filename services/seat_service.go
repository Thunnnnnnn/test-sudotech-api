package services

import (
	"gin-api/models"

	"gin-api/repositories"
)

func GetSeats() ([]models.Seat, error) {
	return repositories.FindAllSeats()
}

func CreateSeat(seat models.Seat) (models.Seat, error) {
	return repositories.CreateSeat(seat)
}

func BookSeat(seatID string) error {
	return repositories.BookSeat(seatID)
}
