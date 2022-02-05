package controllers

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"hackathon-api/configs"
	"hackathon-api/models"
	"hackathon-api/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

//https://www.mongodb.com/blog/post/quick-start-golang--mongodb--data-aggregation-pipeline

var statsCollection = configs.GetCollection(configs.DB, "donations")

func SumDonationsByMoney() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		//money := c.Param("money")
		defer cancel()

		//matchStage := bson.D{{"$match", bson.D{{"moneyType", money}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", "$moneyType"}, {"total", bson.D{{"$sum", "$amount"}}}}}}

		resultCursor, err := statsCollection.Aggregate(ctx, mongo.Pipeline{ /*matchStage,*/ groupStage})
		if err != nil {
			println(err)
		}

		var resultData []models.Statistics
		resultCursor.All(ctx, resultData)

		if err = resultCursor.All(ctx, &resultData); err != nil {
			panic(err)
		}

		var total float64 = 0
		for i, stat := range resultData {
			resultData[i].TotalAmount = stat.Total * models.GetMoney()[stat.Money]
			total = total + resultData[i].TotalAmount
		}

		type StatResponse struct {
			Stats []models.Statistics `json:"stats,omitempty"`
			Total float64             `json:"total,omitempty"`
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, StatResponse{
			Stats: resultData,
			Total: total,
		})

	}

}
