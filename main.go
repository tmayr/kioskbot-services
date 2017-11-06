package main

import (
	"fmt"
	KioskbotLib "kioskbot-services/lib"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv"
	"github.com/robfig/cron"
)

func KioskbotAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, keyExists := os.LookupEnv("KB_SERVICES_API_KEY")

		if !keyExists {
			c.JSON(500, map[string]string{"error": "API_KEY not set in Application"})
			c.Abort()
			return
		}

		if c.GetHeader("X-KB-SERVICES-API-KEY") != apiKey {
			c.JSON(401, map[string]string{"error": "not authorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	// starting the cronjob
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		fmt.Println("Running Email Cron...")
		KioskbotLib.Email()
	})
	c.Start()

	router := gin.New()

	// Middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	{
		v1.Use(KioskbotAuth())
		v1.GET("/update-algolia", func(c *gin.Context) {
			kioskitems := KioskbotLib.FetchProductsFromMongo()
			KioskbotLib.SendProductsToAlgolia(kioskitems)
			c.JSON(http.StatusOK, kioskitems)
		})

		v1.GET("/slack-request", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		v1.GET("/", func(c *gin.Context) {
			endpoints := []string{"/api/v1/update-algolia"}
			c.JSON(200, map[string]interface{}{"endpoints": endpoints})
		})
	}

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Nothing to see here.")
	})

	router.Run()
}
