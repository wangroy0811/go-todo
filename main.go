package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ichtrojan/go-todo/models"
	"github.com/ichtrojan/go-todo/routes"
	"github.com/ichtrojan/thoth"
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger, _ := thoth.Init("log")

	if err := godotenv.Load(); err != nil {
		logger.Log(errors.New("no .env file found"))
		log.Fatal("No .env file found")
	}

	port, exist := os.LookupEnv("PORT")

	if !exist {
		logger.Log(errors.New("PORT not set in .env"))
		log.Fatal("PORT not set in .env")
	}

	// Initialize blog table before starting the server
	if err := models.InitBlogTable(); err != nil {
		logger.Log(errors.New("failed to initialize blog table"))
		log.Fatal("Failed to initialize blog table:", err)
	}

	// Check if port is already in use
	if isPortInUse(port) {
		err := fmt.Errorf("port %s is already in use, please close the application using this port or use a different port", port)
		logger.Log(err)
		log.Fatal(err)
	}

	// Initialize the server with routes
	router := routes.Init()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine so we can gracefully handle shutdown
	go func() {
		log.Printf("Server starting on port %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log(err)
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Log(fmt.Errorf("server shutdown failed: %v", err))
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}

// isPortInUse checks if a port is already in use
func isPortInUse(port string) bool {
	conn, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return true
	}
	conn.Close()
	return false
}