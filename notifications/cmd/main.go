package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/BoaglioFranco/casita/notifications/internal/auth"
	"github.com/BoaglioFranco/casita/notifications/internal/database"
	"github.com/BoaglioFranco/casita/notifications/internal/migrations"
	"github.com/BoaglioFranco/casita/notifications/internal/subscriptions"
	"github.com/gin-gonic/gin"
)

func main() {
	requireAuth, err := auth.NewMiddleware(os.Getenv("S2S_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	db, err := database.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := migrations.Up(databaseURL); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	service := subscriptions.NewNotificationService(db, os.Getenv("VAPID_PUBLIC_KEY"), os.Getenv("VAPID_PRIVATE_KEY"))
	service.SetupPublicNotificationRoutes(router)
	internalRoutes := router.Group("/")
	internalRoutes.Use(requireAuth)
	service.SetupNotificationRoutes(internalRoutes)

	port := os.Getenv("PORT")
	if port == "" {
		port = "2346"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start notifications server: %v", err)
	}
}
