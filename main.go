package main

import (
	"log"
	"net/http"

	"gin-api/config"
	"gin-api/database"
	"gin-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// init gin
	r := gin.Default()

	// init redis
	database.InitRedis()
	if err := database.RDB.Set(database.Ctx, "test", "hello redis cloud", 0).Err(); err != nil {
		log.Fatal("Redis error:", err)
	}

	// connect mongo
	if err := database.ConnectMongo(); err != nil {
		log.Fatal("Mongo connection failed:", err)
	}

	if err := config.GoogleOAuthInit(); err != nil {
		log.Fatal("Google OAuth initialization failed:", err)
	}

	// health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Gin 🚀",
		})
	})

	// routes
	routes.UserRoutes(r)
	routes.SeatRoutes(r)
	routes.AuthRoutes(r)

	// run server
	r.Run(":8080")
}
