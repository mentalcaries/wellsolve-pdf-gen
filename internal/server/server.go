package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port              int
	gotenbergEndpoint string
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	gotenbergEndpoint := os.Getenv("PDF_SERVICE_API")
	if gotenbergEndpoint == "" {
		log.Fatal("PDF_SERVICE_API var must be set")
	}

	NewServer := &Server{
		port:              port,
		gotenbergEndpoint: gotenbergEndpoint,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("PDF service running on [%d]", port)
	return server
}
