package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BoaglioFranco/groceriesList/notifications/internal/migrations"
	"github.com/gin-gonic/gin"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if err := migrations.Up(databaseURL); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "2346"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start notifications server: %v", err)
	}
}
