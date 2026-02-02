package services

import (
	"gin-api/models"

	user_repositories "gin-api/repositories"
)

func GetUsers() ([]models.User, error) {
	return user_repositories.FindAllUsers()
}
