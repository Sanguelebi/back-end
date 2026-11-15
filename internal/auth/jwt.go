package auth

import (
	"time"

	"course-selection-platform/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userData models.LoginToken, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userData.ID,
		"email": userData.Email,
		"role":  userData.Role,
		"name":  userData.FirstName + " " + userData.LastName,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Fixed to match 24 hours
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

