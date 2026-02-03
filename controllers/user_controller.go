package controllers

import (
	"github.com/gin-gonic/gin"

	"gin-api/services"
)

func GetUsers(c *gin.Context) {
	users, err := services.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}

// func GetUserByEmail(c *gin.Context) {
// 	email := c.Query("email")
// 	user, err := services.GetUserByEmail(email)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(200, user)
// }

// func CreateUser(c *gin.Context) {
// 	user, err := services.CreateUser(c)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(201, user)
// }
