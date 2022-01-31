package controllers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hackathon-api/configs"
	"hackathon-api/models"
	"hackathon-api/responses"
	"hackathon-api/services"
	"hackathon-api/utils"
	"net/http"
	"strconv"
	"time"

	pdfgen "github.com/ca-gip/hackathon-reward/pkg/generator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var donationCollection *mongo.Collection = configs.GetCollection(configs.DB, "donations")
var validate = validator.New()

func CreateDonation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var donation models.Donation
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&donation); err != nil {
			c.JSON(http.StatusBadRequest, responses.DonationResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&donation); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.DonationResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		// Validate MoneyType
		if err := utils.ValidateMoneyType(donation.MoneyType); err != nil {
			c.JSON(http.StatusBadRequest, responses.DonationResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Create a new Donation object
		newDonation := models.Donation{
			ID:        primitive.NewObjectID(),
			DonorName: donation.DonorName,
			Amount:    donation.Amount,
			MoneyType: donation.MoneyType,
		}

		hash := uuid.New().String()
		pdfname := fmt.Sprintf("%s.pdf", hash)

		// Generate a pdf document
		pdffile, err := pdfgen.GeneratePerfectDocument(donation.DonorName, donation.Amount, hash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		err = services.UploadFile(pdffile, pdfname)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		newDonation.PDFRef = pdfname

		result, err := donationCollection.InsertOne(ctx, newDonation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.DonationResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetADonation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		donationId := c.Param("donationId")
		var donation models.Donation
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(donationId)

		err := donationCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&donation)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		pdfFile, err := services.DownloadFile(donation.PDFRef)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		donation.PDFfile = pdfFile

		c.JSON(http.StatusOK, responses.DonationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": donation}})
	}
}

func DeleteADonation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		DonationId := c.Param("donationId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(DonationId)

		result, err := donationCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.DonationResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Donation with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.DonationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Donation successfully deleted!"}},
		)
	}
}

func GetAllDonations() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := c.DefaultQuery("limit", "100")

		findOptions := options.Find()

		// Sort by `_id` field descending
		findOptions.SetSort(bson.D{{"_id", -1}})

		// Validate limit is int
		if limitInt, err := strconv.Atoi(limit); err != nil {
			c.JSON(http.StatusBadRequest, responses.DonationResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		} else {
			// Limit number of document returned
			findOptions.SetLimit(int64(limitInt))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var Donations []models.Donation
		defer cancel()

		results, err := donationCollection.Find(ctx, bson.M{}, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleDonation models.Donation
			if err = results.Decode(&singleDonation); err != nil {
				c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			Donations = append(Donations, singleDonation)
		}

		c.JSON(http.StatusOK,
			responses.DonationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": Donations}},
		)
	}
}

func DownloadDonation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		donationId := c.Param("donationId")
		var donation models.Donation
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(donationId)

		err := donationCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&donation)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		pdfFile, err := services.DownloadFile(donation.PDFRef)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		donation.PDFfile = pdfFile

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename=pdf")
		c.Header("Content-Type", "application/octet-stream")
		c.Data(http.StatusOK, "application/octet-stream", pdfFile)

	}
}

func GetMoney() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK,
			responses.DonationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": models.GetMoney()}},
		)
	}
}
