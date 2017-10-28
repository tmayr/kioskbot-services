package routes

import (
	KioskbotLib "kioskbot-services/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Initialize(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		kioskitems := KioskbotLib.FetchProductsFromMongo()
		KioskbotLib.SendProductsToAlgolia(kioskitems)
		c.JSON(http.StatusOK, kioskitems)
	})
}
