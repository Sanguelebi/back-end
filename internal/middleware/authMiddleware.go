package middleware

import (
	"fmt"
	"net/http"

	"course-selection-platform/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		// Try to parse the JWT.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Set key/value pairs from the decoded token to the Gin Request, so that other endpoints don't have to decode it.
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if subFloat, ok := claims["sub"].(float64); ok {
				c.Set("userID", uint(subFloat))
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format in token"})
				return
			}

			c.Set("userRole", claims["role"])
			c.Set("userEmail", claims["email"])
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			return
		}

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.Keys["userRole"]
		if userRole != models.RoleAdmin.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed to access the endpoint."})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ManagerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.Keys["userRole"]
		if userRole != models.RoleManager.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed to access the endpoint."})
			c.Abort()
			return
		}
		c.Next()
	}
}
