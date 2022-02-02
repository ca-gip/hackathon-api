package models

type Statistics struct {
	Money       string  `json:"money,omitempty" bson:"_id,omitempty"`
	Total       float64 `json:"total,omitempty" bson:"total,omitempty"`
	TotalAmount float64 `json:"totalAmount,omitempty"`
}
