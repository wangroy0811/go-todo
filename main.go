package main

import (
	"errors"
	"github.com/ichtrojan/go-todo/models"
	"github.com/ichtrojan/go-todo/routes"
	"github.com/ichtrojan/thoth"
	"log"
	"net/http"
	"os"
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

	err := http.ListenAndServe(":"+port, routes.Init())

	if err != nil {
		logger.Log(err)
		log.Fatal(err)
	}
}
