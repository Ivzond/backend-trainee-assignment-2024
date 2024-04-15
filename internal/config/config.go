package config

import (
	"fmt"
)

const (
	DBHost     = "localhost"
	DBPort     = "5432"
	DBUser     = "postgres"
	DBPassword = "12345678"
	DBName     = "banner_db"
	RedisHost  = "localhost"
	RedisPort  = "6379"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
}

func LoadConfig() (*Config, error) {
	DatabaseURL := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		DBHost,
		DBPort,
		DBUser,
		DBName,
		DBPassword,
	)
	RedisURL := fmt.Sprintf(RedisHost + ":" + RedisPort)

	config := &Config{
		DatabaseURL: DatabaseURL,
		RedisURL:    RedisURL,
	}

	return config, nil
}
