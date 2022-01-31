package controllers

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"hackathon-api/configs"
	"hackathon-api/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

//https://www.mongodb.com/blog/post/quick-start-golang--mongodb--data-aggregation-pipeline

var statsCollection = configs.GetCollection(configs.DB, "donations")

func CountDonationByMoney() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		money := c.Param("money")
		defer cancel()

		count, err := donationCollection.CountDocuments(ctx, bson.M{"moneyType": money})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK,
			responses.DonationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": count}},
		)

	}
}

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

		var resultData []bson.M
		resultCursor.All(ctx, resultData)

		if err = resultCursor.All(ctx, &resultData); err != nil {
			panic(err)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, resultData)

	}

}

func SumDonations() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		groupStage := bson.D{{"$group", bson.D{{"_id", ""}, {"total", bson.D{{"$sum", "$amount"}}}}}}

		resultCursor, err := statsCollection.Aggregate(ctx, mongo.Pipeline{ /*matchStage,*/ groupStage})
		if err != nil {
			println(err)
		}

		var resultData []bson.M
		resultCursor.All(ctx, resultData)

		if err = resultCursor.All(ctx, &resultData); err != nil {
			panic(err)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, resultData)

	}

}
