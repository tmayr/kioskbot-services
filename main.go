package main

import (
	"kioskbot-services/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv"
)

func main() {
	router := gin.New()

	// Middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Routes
	routes.Initialize(router)

	router.Run()
}
