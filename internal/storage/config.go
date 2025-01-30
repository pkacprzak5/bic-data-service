package storage

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type PostgresConfig struct {
	Host     string
	Port     string
	DB_Port  string
	User     string
	Password string
	Database string
}

var Envs = initConfig()

func initConfig() PostgresConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	} else {
		err := godotenv.Overload()
		if err != nil {
			log.Println("Failed to load .env file")
		}
	}
	return PostgresConfig{
		Port:     getEnv("PORT", "8080"),
		DB_Port:  getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "example_user"),
		Password: getEnv("DB_PASSWORD", "Passwd@1234"),
		Host:     getEnv("DB_HOST", "localhost"),
		Database: getEnv("DB_NAME", "bicdatabase"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
