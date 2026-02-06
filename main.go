package main

import (
	"log"
	"net/http"
	"time"

	"gin-api/config"
	"gin-api/database"
	"gin-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	// load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// init gin
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Bearer-Token",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("connected:", s.ID())
		return nil
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal(err)
		}
	}()
	defer server.Close()

	r.GET("/socket.io/*any", gin.WrapH(server))
	r.POST("/socket.io/*any", gin.WrapH(server))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Gin 🚀",
		})
	})
	// routes
	routes.UserRoutes(r)
	routes.SeatRoutes(r)
	routes.AuthRoutes(r)
	routes.TheaterRoutes(r)

	// run server
	r.Run(":8080")
}
