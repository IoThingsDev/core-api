package models

type Card struct {
	Token string `json:"token" bson:"token" binding:"required"`
}
