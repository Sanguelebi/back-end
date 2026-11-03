package services

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/auth"
	"course-selection-platform/internal/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(email, plainPassword string, app *config.AppData) (string, models.LoginUserResponse, error) {
    var user models.Users

    err := app.DB.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return "", models.LoginUserResponse{}, errors.New("invalid email or password")
        }
        return "", models.LoginUserResponse{}, err
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

    loginDTO := models.LoginUserResponse{
        ID:        user.ID,
        Email:     user.Email,
        Role:      user.Role.String(),
        FirstName: user.FirstName,
        LastName:  user.LastName,
    }

    // Generate the token
    tokenString, err := auth.GenerateToken(loginDTO, app.JWTSecret)
    if err != nil {
        return "", models.LoginUserResponse{}, err
    }

    // Return the string token, the response DTO, and nil for error
    return tokenString, userResponse, nil
}