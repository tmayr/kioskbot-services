package main

import (
	"encoding/json"
	KioskbotLib "kioskbot-services/lib"
	KioskTypes "kioskbot-services/types"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv"
	"github.com/nlopes/slack"
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
	KioskbotLib.Email()
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

	router.POST("/api/v1/slack-response", func(c *gin.Context) {
		var jsonPayload slack.AttachmentActionCallback
		formPayload := c.PostForm("payload")

		// payload comes wrapped in ""
		formPayload = strings.TrimSuffix(formPayload, "\"")

		// string is still url encoded
		formPayload, err := url.QueryUnescape(formPayload)
		if err != nil {
			panic(err)
		}

		bytes := []byte(formPayload)
		err = json.Unmarshal(bytes, &jsonPayload)
		if err != nil {
			panic(err)
		}

		// as soon as we have the serialized json, check if its an authorized one
		if jsonPayload.Token != os.Getenv("SLACK_APP_TOKEN") {
			c.JSON(400, map[string]string{"error": "invald slack app token"})
			return
		}

		callbackParts := strings.Split(jsonPayload.CallbackID, "?")
		callbackID := callbackParts[0]
		callbackQueryValues, _ := url.ParseQuery(callbackParts[1])

		if callbackID == "wire_user_selection" {
			KioskbotLib.SendPayment(KioskTypes.Payment{
				User:   jsonPayload.Actions[0].SelectedOptions[0].Value,
				Amount: callbackQueryValues["amount"][0],
			})
		}

		c.JSON(200, jsonPayload)
	})

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Nothing to see here.")
	})

	router.Run()
}
