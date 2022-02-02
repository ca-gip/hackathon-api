package routes

import (
	"hackathon-api/controllers"

	"github.com/gin-gonic/gin"
)

func StatisticsRoute(router *gin.Engine) {
	router.GET("/statistics", controllers.SumDonationsByMoney())
}
