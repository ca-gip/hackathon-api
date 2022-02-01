package models

// Stats type
type Stats struct {
	ID    string `json:"_id,omitempty" bson:"_id,omitempty"`
	total string `json:"total,omitempty" bson:"total,omitempty"`
}
