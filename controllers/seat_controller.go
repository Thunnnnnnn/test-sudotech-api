package controllers

import (
	"github.com/gin-gonic/gin"

	seat_services "gin-api/services"
)

func GetSeats(c *gin.Context) {
	seats, err := seat_services.GetSeats()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, seats)
}
