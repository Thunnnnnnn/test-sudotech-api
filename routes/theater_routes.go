package routes

import (
	"gin-api/controllers"

	"github.com/gin-gonic/gin"
)

func TheaterRoutes(r *gin.Engine) {
	theaters := r.Group("/theaters")
	{
		theaters.GET("", controllers.GetTheaters)
		theaters.POST("", controllers.CreateTheater)
	}
}
