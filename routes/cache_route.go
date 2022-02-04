package routes

import (
	"hackathon-api/controllers"

	"github.com/gin-gonic/gin"
)

func CacheRoute(router *gin.Engine) {
	router.GET("/cache/items", controllers.CacheContent())
	router.GET("/cache/keys", controllers.CacheKeys())
}
