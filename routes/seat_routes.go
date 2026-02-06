package routes

import (
	"gin-api/controllers"
	"gin-api/middleware"

	"github.com/gin-gonic/gin"
)

func SeatRoutes(r *gin.Engine) {
	seats := r.Group("/seats")
	seats.Use(middleware.AuthMiddleware())
	{
		seats.GET("", controllers.GetSeats)
		seats.POST("", controllers.CreateSeat)
		seats.POST("/book/:id", controllers.BookSeat)
		seats.POST("/cancel/:id", controllers.CancelSeatBooking)
		seats.POST("/confirm/:id", controllers.ConfirmSeatBooking)
	}
}
