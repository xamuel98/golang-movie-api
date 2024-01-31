package main

import (
	"fmt"
	"log"
	routes "movie-api/api/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get PORT from .env
	port := os.Getenv("PORT")

	// If PORT is not declared, set to 0.0.0.0:8080
	if port == "" {
		port = "8080"
	}

	// Initialize Gin router
	router := gin.Default()
	router.Use(gin.Logger())

	// Use the routes
	routes.MoviesRoutes(router)
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	fmt.Printf("Starting server on port: %v\n", port)

	// Start the server
	router.Run(":" + port) // listen and serve on port
}
