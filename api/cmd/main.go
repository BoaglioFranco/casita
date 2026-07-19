package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/BoaglioFranco/groceriesList/api/internal/auth"
	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/BoaglioFranco/groceriesList/api/internal/households"
	"github.com/BoaglioFranco/groceriesList/api/internal/items"
	"github.com/BoaglioFranco/groceriesList/api/internal/lists"
	"github.com/BoaglioFranco/groceriesList/api/internal/notifications"
	"github.com/gin-gonic/gin"
)

func main() {
	vapidPrivateKey := os.Getenv("VAPID_PRIVATE_KEY")
	vapidPublicKey := os.Getenv("VAPID_PUBLIC_KEY")

	router := gin.Default()

	// Quick demo CORS handling for the static local UI.
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "groceries.db"
	}

	db, err := database.New(context.Background(), dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	authService := auth.NewAuthService(db)
	notificationsService := notifications.NewNotificationService(db, vapidPublicKey, vapidPrivateKey)
	listsService := lists.NewListService(db)
	itemsService := items.NewItemService(db, notificationsService)
	householdService := households.NewHouseholdService(db)

	authService.SetupAuthRoutes(router)
	notificationsService.SetupPublicNotificationRoutes(router)

	authenticatedRoutes := router.Group("/")
	authenticatedRoutes.Use(authService.RequireAuth)

	householdRoutes := authenticatedRoutes.Group("/households/:householdId")
	householdRoutes.Use(householdService.RequireMembership)
	listsService.SetupListRoutes(householdRoutes)
	itemsService.SetupItemRoutes(householdRoutes)
	notificationsService.SetupNotificationRoutes(householdRoutes)

	port := os.Getenv("PORT")
	if port == "" {
		port = "2345"
	}

	router.Run(":" + port)
}
