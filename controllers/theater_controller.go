package controllers

import (
	"github.com/gin-gonic/gin"

	"gin-api/models"
	"gin-api/services"
)

func GetTheaters(c *gin.Context) {
	theaters, err := services.GetTheaters()
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "data": theaters})
}

func CreateTheater(c *gin.Context) {
	var theater models.Theater

	if err := c.ShouldBindJSON(&theater); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	_, err := services.CreateTheater(theater)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"code": 201, "message": "สร้างที่โรงภาพยนตร์สำเร็จ"})
}

func GetTheater(c *gin.Context) {
	theaterID := c.Param("id")
	theater, err := services.GetTheaterByID(theaterID)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "data": theater})
}
