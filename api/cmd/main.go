package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/BoaglioFranco/groceriesList/api/internal/items"
	"github.com/BoaglioFranco/groceriesList/api/internal/lists"
	"github.com/BoaglioFranco/groceriesList/api/internal/middleware"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/gin-gonic/gin"
)

func generateVapidKeys() (string, string, error) {
	newPrivateKey, newPublicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatalf("Failed to generate VAPID keys: %v", err)
	}
	return newPrivateKey, newPublicKey, nil
}

func demo(router *gin.Engine, vapidPublicKey string, vapidPrivateKey string, currentSubscription *webpush.Subscription) {

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"response": "world",
		})
	})

	router.GET("/vapid-public-key", func(c *gin.Context) {
		if vapidPublicKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "VAPID public key not set",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"publicKey": vapidPublicKey,
		})
	})

	router.POST("/subscribe", func(c *gin.Context) {
		var subscription webpush.Subscription
		if err := c.ShouldBindJSON(&subscription); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription data"})
			return
		}
		currentSubscription = &subscription
		c.Status(http.StatusOK)
	})

	router.POST("/send-test", func(c *gin.Context) {

		if currentSubscription == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no subscription yet"})
			return
		}

		response, err := webpush.SendNotification(
			[]byte(`{"title":"Groceries","body":"Hello from Go"}`),
			currentSubscription,
			&webpush.Options{
				Subscriber:      "mailto:test@example.com",
				VAPIDPublicKey:  vapidPublicKey,
				VAPIDPrivateKey: vapidPrivateKey,
				TTL:             60,
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer response.Body.Close()

		c.Status(http.StatusOK)
	})
}

func main() {
	vapidPrivateKey := os.Getenv("VAPID_PRIVATE_KEY")
	vapidPublicKey := os.Getenv("VAPID_PUBLIC_KEY")
	var currentSubscription *webpush.Subscription

	router := gin.Default()

	// Quick demo CORS handling for the static local UI.
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	demo(router, vapidPublicKey, vapidPrivateKey, currentSubscription)

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "groceries.db"
	}

	db, err := database.New(context.Background(), dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	listsService := lists.NewListService(db)
	itemsService := items.NewItemService(db)

	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.RequireMember())
	listsService.SetupListRoutes(protectedRoutes)
	itemsService.SetupItemRoutes(protectedRoutes)

	port := os.Getenv("PORT")
	if port == "" {
		port = "2345"
	}

	router.Run(":" + port)
}
