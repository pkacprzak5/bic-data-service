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
		Port:     GetEnv("PORT", "8080"),
		DB_Port:  GetEnv("DB_PORT", "5432"),
		User:     GetEnv("DB_USER", "example_user"),
		Password: GetEnv("DB_PASSWORD", "Passwd@1234"),
		Host:     GetEnv("DB_HOST", "localhost"),
		Database: GetEnv("DB_NAME", "bicdatabase"),
	}
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
