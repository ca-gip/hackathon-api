package routes

import (
	"hackathon-api/controllers"

	"github.com/gin-gonic/gin"
)

func StatisticsRoute(router *gin.Engine) {
	router.GET("/statistics/:money", controllers.CountDonationByMoney())
	router.GET("/statistics/sum", controllers.SumDonationsByMoney())
	router.GET("/statistics/total", controllers.SumDonations())
}
