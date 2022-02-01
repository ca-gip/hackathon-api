package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"hackathon-api/configs"
	"hackathon-api/routes"
	"time"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"OPTIONS", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//run database
	configs.ConnectDB()

	//routes
	routes.DonationRoute(router)
	routes.StatisticsRoute(router)

	router.Run("0.0.0.0:8080")
}
