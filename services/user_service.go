package services

import (
	"gin-api/helpers"
	"gin-api/models"
	"gin-api/repositories"
)

func GetUsers() ([]models.User, error) {
	return repositories.FindAllUsers()
}

func FindOrCreateGoogleUser(google map[string]interface{}) (models.User, error) {
	email := helpers.GetString(google, "email")

	user, err := repositories.FindUserByEmail(email)
	if err != nil {
		panic("err:" + err.Error())
		// return models.User{}, err
	}

	if !user.ID.IsZero() {
		return user, nil
	}

	newUser := models.User{
		Email:         email,
		Firstname:     helpers.GetString(google, "given_name"),
		Surname:       helpers.GetString(google, "family_name"),
		FullName:      helpers.GetString(google, "name"),
		Picture:       helpers.GetString(google, "picture"),
		VerifiedEmail: helpers.GetBool(google, "verified_email"),
	}

	return repositories.CreateUser(newUser)
}
