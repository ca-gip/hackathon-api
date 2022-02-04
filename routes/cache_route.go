package routes

import (
	"hackathon-api/controllers"

	"github.com/gin-gonic/gin"
)

func CacheRoute(router *gin.Engine) {
	router.GET("/cache", controllers.CacheContent())
}
