package main

import (
	"log"
	"net/http"

	"gin-api/database"
	redisclient "gin-api/redis"
	routes "gin-api/routes"

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
	rdb := redisclient.NewRedisClient()
	if err := rdb.Set(redisclient.Ctx, "test", "hello redis cloud", 0).Err(); err != nil {
		log.Fatal("Redis error:", err)
	}

	// connect mongo
	if err := database.ConnectMongo(); err != nil {
		log.Fatal("Mongo connection failed:", err)
	}

	// health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Gin 🚀",
		})
	})

	// routes
	routes.UserRoutes(r)

	// run server
	r.Run(":8080")
}
