package main

import (
	"context"
	"fmt"
	"github.com/pkacprzak5/bic-data-service/internal/app"
	"github.com/pkacprzak5/bic-data-service/internal/storage"
	"log"
	"os"
	"os/signal"
)

func main() {
	dbConfig := storage.PostgresConfig{
		Host:     storage.Envs.Host,
		DB_Port:  storage.Envs.DB_Port,
		User:     storage.Envs.User,
		Password: storage.Envs.Password,
		Database: storage.Envs.Database,
	}
	postgresDB, err := storage.NewPostgreSQLStorage(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL storage: %v", err)
	}

	db, err := postgresDB.Init()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	port := fmt.Sprintf(":%v", storage.Envs.Port)

	store := storage.NewRelationalDB(db)
	api := app.NewAPIServer(port, store)
	err = api.Start(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
}
