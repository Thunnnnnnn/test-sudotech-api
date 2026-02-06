package controllers

import (
	"github.com/gin-gonic/gin"

	"gin-api/models"
	"gin-api/services"
)

func GetSeats(c *gin.Context) {
	seats, err := services.GetSeats()
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "data": seats})
}

func CreateSeat(c *gin.Context) {
	var seat models.Seat

	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	_, err := services.CreateSeat(seat)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"code": 201, "message": "สร้างที่นั่งสำเร็จ"})
}

func BookSeat(c *gin.Context) {
	seatID := c.Param("id")
	err := services.BookSeat(c, seatID)
	if err != nil {
		if err.Error() == "มีคนกำลังจองอยู่" {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "มีคนกำลังจองอยู่",
			})
			return
		}

		if err.Error() == "คุณจองเก้าอี้นี้อยู่" {
			c.JSON(200, gin.H{
				"code":    200,
				"message": "คุณจองเก้าอี้นี้อยู่",
			})
			return
		}
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "จองที่นั่งสำเร็จ"})
}

func CancelSeatBooking(c *gin.Context) {
	seatID := c.Param("id")
	err := services.CancelSeatBooking(c, seatID)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "message": "ยกเลิกการจองที่นั่งสำเร็จ"})
}

func ConfirmSeatBooking(c *gin.Context) {
	seatID := c.Param("id")
	err := services.ConfirmSeatBooking(c, seatID)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "message": "ยืนยันการจองที่นั่งสำเร็จ"})
}
