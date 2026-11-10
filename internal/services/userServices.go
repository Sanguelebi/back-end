package services

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/auth"
	"course-selection-platform/internal/models"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginService(email, plainPassword string, app *config.AppData) (string, models.LoginUserResponse, error) {
	var user models.Users

	err := app.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", models.LoginUserResponse{}, errors.New("invalid email or password")
		}
		return "", models.LoginUserResponse{}, err
	}

	if len(user.Password) == 0 {
		return "", models.LoginUserResponse{}, errors.New("invalid email or password")
	}

	if len(plainPassword) > 72 {
		return "", models.LoginUserResponse{}, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword))
	if err != nil {
		return "", models.LoginUserResponse{}, errors.New("invalid email or password")
	}

	userResponse := models.LoginUserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Role:          user.Role.String(),
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		PhotoUrl:      user.PhotoUrl,
		StudentStatus: user.StudentStatus,
	}

	loginDTO := userResponse.CreateLoginToken()

	// Generate the token
	tokenString, err := auth.GenerateToken(*loginDTO, app.JWTSecret)
	if err != nil {
		return "", models.LoginUserResponse{}, err
	}

	// Return the string token, the response DTO, and nil for error
	return tokenString, userResponse, nil
}

func GetAllUsersService(app *config.AppData) ([]models.Users, error) {
	var users []models.Users
	result := app.DB.Find(&users)
	return users, result.Error
}

func AddUserService(userData *models.RegisterUserRequest, app *config.AppData) (*uint, int, error) {
	var count int64
	err := app.DB.Model(&models.Users{}).
		Where("email = ?", userData.Email).
		Count(&count).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if count > 0 {
		return nil, http.StatusBadRequest, errors.New("User with this E-Mail already exists.")
	}

	if userData.StudentStatus == nil && userData.Role == models.RoleStudent {
		return nil, http.StatusBadRequest, errors.New("Specify Student Status for the student.")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	userData.Password = string(hashedPassword)

	// Create the new user.
	user := userData.CreateTableUser()
	result := app.DB.Model(&models.Users{}).Create(user)
	if result.Error != nil {
		return nil, http.StatusInternalServerError, result.Error
	}
	return &user.ID, http.StatusCreated, nil
}
