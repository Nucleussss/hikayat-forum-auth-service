package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost    string `env:"DB_HOST"`
	DBUser    string `env:"DB_USER"`
	DBPass    string `env:"DB_PASSWORD"`
	DBName    string `env:"DB_NAME"`
	DBPort    string `env:"DB_PORT"`
	DBsslMode string `env:"DB_SSL_MODE"`

	GRPCPort string `env:"GRPC_PORT"`
}

func LoadConfig() *Config {
	// load configurations from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v ", err)
		return nil
	}

	return &Config{
		DBHost:    os.Getenv("DB_HOST"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASSWORD"),
		DBName:    os.Getenv("DB_NAME"),
		DBPort:    os.Getenv("DB_PORT"),
		DBsslMode: os.Getenv("DB_SSL_MODE"),

		GRPCPort: os.Getenv("GRPC_PORT"),
	}
}
