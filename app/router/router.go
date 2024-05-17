package router

import (
	"currency_mail/app/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/rate", handlers.RateHandler)

	// r.GET("/subscribes", handlers.RateHandler)
	r.POST("/subscribe", handlers.SubscribeHandler)
	r.POST("/unsubscribe", handlers.UnsubscribeHandler)
	// r.POST("/update_subscribes", handlers.SubscribeHandler)

	return r
}
