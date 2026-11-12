package handlers

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/models"
	"course-selection-platform/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func LoginHandler(app *config.AppData) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginUserRequest

		// ShouldBindJSON deserializes an HTTP Request into JSON body of a Go struct.
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		tokenString, userData, err := services.LoginService(req.Email, req.Password, app)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Set the cookie
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

func GetAllUsersHandler(app *config.AppData) gin.HandlerFunc {
	return func(c *gin.Context) {
		userData, err := services.GetAllUsersService(app)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get users from DB."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Fetched users succesfully", "users": userData})
	}
}

func AddUserHandler(app *config.AppData) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegisterUserRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		newUserId, status, err := services.AddUserService(&req, app)
		if err != nil {
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}

		c.JSON(status, gin.H{"message": "User created succesfully!", "id": newUserId})
	}
}

func EditUserHandler(app *config.AppData) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.EditUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		idParam := c.Param("id")
		parsedID, err := strconv.ParseUint(idParam, 10, 0)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID Format."})
			return
		}

		status, err := services.EditUserService(uint(parsedID), &req, app)
		if err != nil {
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}

		c.JSON(status, gin.H{"message": "User edited successfully!"})
	}
}
