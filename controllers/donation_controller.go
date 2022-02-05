package controllers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hackathon-api/configs"
	"hackathon-api/models"
	"hackathon-api/responses"
	"hackathon-api/services"
	"hackathon-api/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os/exec"
)

var donationCollection = configs.GetCollection(configs.DB, "donations")
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
			c.JSON(http.StatusNotAcceptable, responses.DonationResponse{Status: http.StatusNotAcceptable, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Validate minmum amount
		if donation.Amount <= 0.0001 {
			c.JSON(http.StatusBadRequest, responses.DonationResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "amount must be equals or superior to 0.0001"}})
			return
		}

		// Validate unique constraint on name and money
		pipeline := bson.D{
			{"$and", []interface{}{
				bson.D{{"donatorName", donation.DonorName}},
				bson.D{{"moneyType", donation.MoneyType}},
			}},
		}

		count, err := donationCollection.CountDocuments(ctx, pipeline)

		if err != nil {
			c.JSON(http.StatusConflict, responses.DonationResponse{Status: http.StatusConflict, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if count >= 1 {
			c.JSON(http.StatusConflict, responses.DonationResponse{Status: http.StatusConflict, Message: "error", Data: map[string]interface{}{"data": fmt.Sprintf("donor name with %s currency already existing in database", donation.MoneyType)}})
			return
		}

		// Create a new Donation object
		newDonation := models.Donation{
			ID:        primitive.NewObjectID(),
			DonorName: donation.DonorName,
			Amount:    donation.Amount,
			MoneyType: donation.MoneyType,
			PDFSize:   donation.PDFSize,
		}

		hash := uuid.New().String()

		args := []string{"create", "--output", "Bytes", "--donor", donation.DonorName, "--hash", hash, "--currency", donation.MoneyType, "--amount", fmt.Sprintf("%v", donation.Amount)}
		pdffile, err := exec.Command("hackathon-reward", args...).Output()

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		pdfsize, err := services.UploadFile(pdffile, hash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		newDonation.PDFSize = pdfsize
		newDonation.PDFRef = hash

		result, err := donationCollection.InsertOne(ctx, newDonation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		newDonation.PDFSize = pdfsize

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

func GetAllDonationsPaginated() gin.HandlerFunc {
	return func(c *gin.Context) {
		sortBy := c.DefaultQuery("sortBy", "amount")
		sortDesc, errSortDesc := strconv.ParseInt(c.DefaultQuery("sortDesc", "-1"), 10, 64)
		page, errPageConvert := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
		itemsPerPage, errPerPage := strconv.ParseInt(c.DefaultQuery("itemsPerPage", "10"), 10, 64)
		searchTerm := c.DefaultQuery("term", "")

		if errSortDesc != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid sort direction, must be 1 or -1 not %i", sortDesc))
			return
		}

		if errPageConvert != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid page number direction, must be an integer: %v", page))
			return
		}

		if errPerPage != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid number of item per page, must be an integer: %v", itemsPerPage))
			return
		}

		donations, found := queryCache.Get(c.Request.RequestURI)

		if found {
			c.JSON(http.StatusOK, donations)
			return
		}

		skipItems := itemsPerPage * (page - 1)
		findOptions := options.FindOptions{
			Limit: &itemsPerPage,
			Skip:  &skipItems,
			Sort:  bson.D{{sortBy, sortDesc}},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var Donations = make([]models.Donation, 0)
		defer cancel()

		// No search term
		var pipeline = bson.D{}

		if len(searchTerm) > 0 {
			likeFilter := bson.M{
				"$regex": primitive.Regex{
					Pattern: "^.*" + searchTerm + ".*",
					Options: "i",
				},
			}

			pipeline = bson.D{
				{"$or", []interface{}{
					bson.D{{"donatorName", likeFilter}},
					bson.D{{"pdfRef", likeFilter}},
				}},
			}

		}

		results, err := donationCollection.Find(ctx, pipeline, &findOptions)
		count, err := donationCollection.CountDocuments(ctx, pipeline)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		if err = results.All(ctx, &Donations); err != nil {
			log.Err(err)
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		data := responses.DonationResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": Donations, "total": count}}
		queryCache.Set(c.Request.RequestURI, data, 60*time.Second)
		c.JSON(http.StatusOK, data)
	}
}

func DownloadDonation() gin.HandlerFunc {
	return func(c *gin.Context) {
		pdfRef := c.Param("pdfRef")

		pdfFile, err := services.DownloadFile(pdfRef)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.DonationResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+pdfRef+".pdf")
		c.Header("Content-Type", "application/octet-stream")
		c.Data(http.StatusOK, "application/octet-stream", pdfFile)

	}
}

func GetMoney() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, models.GetMoney())
	}
}
