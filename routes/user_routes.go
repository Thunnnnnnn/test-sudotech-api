package routes

import (
	"gin-api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("", controllers.GetUsers)
	}
}
