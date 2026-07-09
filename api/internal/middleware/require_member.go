package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credentials are required to access this resource"})
			c.Abort()
			return
		}

		memberIDText := strings.TrimSpace(authHeader[len("Bearer "):])
		memberID, err := strconv.ParseInt(memberIDText, 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid member id"})
			c.Abort()
			return
		}

		c.Set("member_id", memberID)
		c.Next()
	}
}
