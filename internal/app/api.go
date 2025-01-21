package app

import (
	"context"
	"github.com/pkacprzak5/bic-data-service/internal/storage"
	"log"
	"net/http"
	"time"
)

type APIServer struct {
	address string
	storage storage.Storage
}

func NewAPIServer(address string, storage storage.Storage) *APIServer {
	return &APIServer{address: address, storage: storage}
}

func (s *APIServer) Start(ctx context.Context) error {
	router := http.NewServeMux()
	subrouter := http.NewServeMux()

	router.Handle("/v1/", http.StripPrefix("/v1/", subrouter))

	bankService := NewBankService(s.storage)
	bankService.RegisterRoutes(subrouter)

	server := &http.Server{
		Addr:    s.address,
		Handler: subrouter,
	}

	ch := make(chan error, 1)
	go func() {
		log.Println("Starting API server at", s.address)

		err := server.ListenAndServe()
		if err != nil {
			ch <- err
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		log.Println("Received shutdown signal, shutting down the server...")
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
