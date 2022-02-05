package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Donation type
type Donation struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	DonorName string             `json:"donatorName,omitempty" bson:"donatorName,omitempty"`
	Amount    float32            `json:"amount,omitempty" bson:"amount,omitempty"`
	MoneyType string             `json:"moneyType,omitempty" bson:"moneyType,omitempty"`
	PDFRef    string             `json:"pdfRef,omitempty" bson:"pdfRef,omitempty"`
	PDFfile   []byte             `json:"pdfFile,omitempty" bson:"pdfFile,omitempty"`
	PDFSize   int                `json:"pdfSize,omitempty" bson:"pdfSize,omitempty"`
}

func GetMoney() map[string]float64 {
	var moneys = make(map[string]float64)

	moneys["BTC"] = 3.85
	moneys["ETH"] = 2.68
	moneys["LTC"] = 1.09
	moneys["XMR"] = 1.45
	moneys["EuR"] = 1.12
	moneys["USD"] = 1

	return moneys
}
