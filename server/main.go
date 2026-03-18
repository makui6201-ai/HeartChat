package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"heartchat-server/handlers"
)

func main() {
	r := gin.Default()

	// CORS – set CORS_ORIGIN env var to restrict origins in production.
	// Example: CORS_ORIGIN=https://servicewechat.com,https://your-domain.com
	allowOrigins := []string{"*"}
	if env := os.Getenv("CORS_ORIGIN"); env != "" {
		allowOrigins = strings.Split(env, ",")
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	api := r.Group("/api")
	{
		api.POST("/chat", handlers.Chat)
		api.POST("/greeting", handlers.Greeting)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("HeartChat server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
