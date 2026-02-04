package services

import (
	"gin-api/models"
	"gin-api/repositories"
)

func GetTheaters() ([]models.Theater, error) {
	return repositories.FindAllTheaters()
}

func CreateTheater(theater models.Theater) (models.Theater, error) {
	return repositories.CreateTheater(theater)
}
