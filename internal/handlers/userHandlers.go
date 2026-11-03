package handlers

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/models"
	"course-selection-platform/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(app *config.AppData) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginUserRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		tokenString, userData, err := services.LoginService(req.Email, req.Password, app)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// 3. Set the cookie
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(
			"token",
			tokenString, 
			86400,       
			"/",         
			"",          
			app.IsSecure,
			true,        
		)

		userResponse := models.LoginUserResponse{
			ID:            userData.ID,
			Email:         userData.Email,
			Role:          userData.Role,
			FirstName:     userData.FirstName,
			LastName:      userData.LastName,
			PhotoUrl:      userData.PhotoUrl,
			StudentStatus: userData.StudentStatus,
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": userResponse})
	}
}