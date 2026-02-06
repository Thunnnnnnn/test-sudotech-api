package routes

import (
	"gin-api/controllers"
	"gin-api/middleware"

	"github.com/gin-gonic/gin"
)

func TheaterRoutes(r *gin.Engine) {
	theaters := r.Group("/theaters")
	theaters.Use(middleware.AuthMiddleware())
	{
		theaters.GET("", controllers.GetTheaters)
		theaters.GET("/:id", controllers.GetTheater)
		theaters.POST("", controllers.CreateTheater)
	}
}
