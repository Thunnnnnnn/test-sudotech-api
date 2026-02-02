package controllers

import (
	"github.com/gin-gonic/gin"

	user_services "gin-api/services"
)

func GetUsers(c *gin.Context) {
	users, err := user_services.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}
