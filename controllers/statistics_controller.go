package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"hackathon-api/configs"
	"hackathon-api/models"
	"hackathon-api/responses"
	"net/http"
	"time"
)

//https://www.mongodb.com/blog/post/quick-start-golang--mongodb--data-aggregation-pipeline

var statsCollection = configs.GetCollection(configs.DB, "donations")

func SumDonationsByMoney() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		stats, found := queryCache.Get("stats")

		if found {
			c.JSON(http.StatusOK, stats)
			return
		}

		groupStage := bson.D{{"$group", bson.D{{"_id", "$moneyType"}, {"total", bson.D{{"$sum", "$amount"}}}}}}
		resultCursor, err := statsCollection.Aggregate(ctx, mongo.Pipeline{ /*matchStage,*/ groupStage})
		count, err := statsCollection.CountDocuments(ctx, bson.M{})

		if err != nil {
			log.Err(err)
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return

		}

		var resultData []models.Statistics

		if err = resultCursor.All(ctx, &resultData); err != nil {
			log.Err(err)
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var total float64 = 0
		for i, stat := range resultData {
			resultData[i].TotalAmount = stat.Total * models.GetMoney()[stat.Money]
			total = total + resultData[i].TotalAmount
		}

		type StatResponse struct {
			Stats []models.Statistics `json:"stats,omitempty"`
			Total float64             `json:"total,omitempty"`
			Count int64               `json:"count,omitempty"`
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		queryCache.Set("stats", StatResponse{
			Stats: resultData,
			Total: total,
			Count: count,
		}, 30*time.Second)

		c.JSON(http.StatusOK, StatResponse{
			Stats: resultData,
			Total: total,
			Count: count,
		})

	}

}
