package services

import (
	"gin-api/models"

	"gin-api/repositories"
)

func GetUsers() ([]models.User, error) {
	return repositories.FindAllUsers()
}
