package config

import (
	"fmt"
	"os"
	"strconv"
)

// TODO: Move DBConfig to be a part of the AppConfig struct.
type AppConfig struct {
	JWTSecret []byte
}

// This will be the main config for DB connection.
// Since this is a config folder, we only want configuration files here.
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadAppConfig() (*AppConfig, error) {
	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		return nil, fmt.Errorf("JWT Secret does not exist.")
	}

	return &AppConfig{
		JWTSecret: []byte(jwt_secret),
	}, nil
}

// LoadDBConfig is the function that will populate the DBConfig and return the object of it.
func LoadDBConfig() (*DBConfig, error) {
	required := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"}
	for _, key := range required {
		if _, ok := os.LookupEnv(key); !ok {
			return nil, fmt.Errorf("missing environment variable: %s", key)
		}
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("DB_PORT must be a valid port number")
	}

	return &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}, nil
}