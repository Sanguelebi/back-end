package auth

import (
	"os"
	"time"

	"course-selection-platform/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userData models.LoginUserDTO) (string, error) {
	claims := jwt.MapClaims{
		"sub": userData.ID,
		"email": userData.Email,
		"role":  userData.Role,
		"name":  userData.FirstName + " " + userData.LastName,
		"exp": time.Now().Add(time.Hour).Unix(), // Expires in 24 hours
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}