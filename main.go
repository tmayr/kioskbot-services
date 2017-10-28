package main

import (
	KioskbotLib "kioskbot-services/lib"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv"
)

func main() {
	router := gin.New()

	// Middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		kioskitems := KioskbotLib.FetchProductsFromMongo()
		KioskbotLib.SendProductsToAlgolia(kioskitems)
		c.JSON(http.StatusOK, kioskitems)
	})

	router.Run()
}
