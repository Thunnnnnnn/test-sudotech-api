package services

import (
	"gin-api/models"

	user_repositories "gin-api/repositories"
)

func GetSeats() ([]models.Seat, error) {
	return user_repositories.FindAllSeats()
}
