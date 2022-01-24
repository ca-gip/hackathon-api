package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Donation type
type Donation struct {
	ID        	primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	DonorName 	string             `json:"donatorName,omitempty" bson:"donatorName,omitempty"`
	Amount    	string             `json:"amount,omitempty" bson:"amount,omitempty"`
	MoneyType   string             `json:"moneyType,omitempty" bson:"moneyType,omitempty"`
	PDFRef      string             `json:"pdfRef,omitempty" bson:"pdfRef,omitempty"`
}
