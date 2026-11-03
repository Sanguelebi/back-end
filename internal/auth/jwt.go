package auth

import (
	"fmt"
	"os"
	"time"

	"course-selection-platform/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userData models.LoginUserDTO) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", fmt.Errorf("JWT_SECRET environment variable is not set")
    }
    jwtKey := []byte(secret)

    claims := jwt.MapClaims{
        "sub":   userData.ID,
        "email": userData.Email,
        "role":  userData.Role,
        "name":  userData.FirstName + " " + userData.LastName,
        "exp":   time.Now().Add(time.Hour * 24).Unix(), // Fixed to match 24 hours
        "iat":   time.Now().Unix(),
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}