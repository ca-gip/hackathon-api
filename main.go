package main

import (
	"hackathon-api/configs"
	"hackathon-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//run database
	configs.ConnectDB()

	//routes
	routes.DonationRoute(router)

	router.Run("0.0.0.0:8080")
}
