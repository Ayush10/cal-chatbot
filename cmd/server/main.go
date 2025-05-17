package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yourusername/cal-chatbot/internal/api"
	"github.com/yourusername/cal-chatbot/internal/chatbot"
)

func main() {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: failed to get current directory: %v", err)
	}

	// Load environment variables from .env file
	envPath := filepath.Join(dir, ".env")
	log.Printf("Loading .env from: %s", envPath)
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: .env file not found at %s, using environment variables: %v", envPath, err)
	}

	// Set up Gin
	if os.Getenv("DEBUG") != "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a new chatbot instance
	bot, err := chatbot.NewChatbot()
	if err != nil {
		log.Fatalf("Failed to create chatbot: %v", err)
	}

	// Create a new router
	router := gin.Default()

	// Setup CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Create API handlers
	handler := api.NewHandler(bot)
	handler.SetupRoutes(router)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Server starting on port %s...", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
