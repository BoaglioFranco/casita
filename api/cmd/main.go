package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/BoaglioFranco/casita/api/internal/auth"
	"github.com/BoaglioFranco/casita/api/internal/calendar"
	"github.com/BoaglioFranco/casita/api/internal/database"
	"github.com/BoaglioFranco/casita/api/internal/items"
	"github.com/BoaglioFranco/casita/api/internal/lists"
	"github.com/BoaglioFranco/casita/api/internal/notifications"
	"github.com/BoaglioFranco/casita/api/migrations"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Quick demo CORS handling for the static local UI.
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

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

	setupRoutes(router, db, os.Getenv("NOTIFICATIONS_SERVICE_URL"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "2345"
	}

	router.Run(":" + port)
}

// setupRoutes is shared by production startup and HTTP integration tests.
func setupRoutes(router gin.IRouter, db *database.Database, notificationServerURL string) {
	authService := auth.NewAuthService(db)
	notificationsService := notifications.NewNotificationService(db, notificationServerURL)
	listsService := lists.NewListService(db)
	itemsService := items.NewItemService(db, notificationsService)
	calendarService := calendar.NewCalendarService(db)

	authService.SetupAuthRoutes(router)
	notificationsService.SetupPublicNotificationRoutes(router)

	authenticatedRoutes := router.Group("/")
	authenticatedRoutes.Use(authService.RequireAuth)
	authService.SetupCurrentUserRoutes(authenticatedRoutes)
	listsService.SetupListRoutes(authenticatedRoutes)
	itemsService.SetupItemRoutes(authenticatedRoutes)
	notificationsService.SetupNotificationRoutes(authenticatedRoutes)
	calendarService.SetupCalendarRoutes(authenticatedRoutes)
}
