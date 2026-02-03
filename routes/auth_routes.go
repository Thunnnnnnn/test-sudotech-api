package routes

import (
	"gin-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.GET("/google", controllers.GoogleLogin)
		auth.GET("/google/callback", controllers.GoogleCallback)
	}
}
