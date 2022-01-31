package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Money struct {
	MoneyType     string
	MoneyUSDPrice float32
}

// Donation type
type Donation struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	DonorName string             `json:"donatorName,omitempty" bson:"donatorName,omitempty"`
	Amount    int                `json:"amount,omitempty" bson:"amount,omitempty"`
	MoneyType string             `json:"moneyType,omitempty" bson:"moneyType,omitempty"`
	PDFRef    string             `json:"pdfRef,omitempty" bson:"pdfRef,omitempty"`
	PDFfile   []byte             `json:"pdfFile,omitempty" bson:"pdfFile,omitempty"`
}

func GetMoney() []Money {
	return []Money{
		{
			"BTC",
			38509,
		},
		{
			"ETH",
			2683,
		},
		{
			"LTC",
			109,
		},
		{
			"XMR",
			145,
		},
		{
			"EUR",
			1.12,
		},
		{
			"USD",
			1,
		},
	}
}
