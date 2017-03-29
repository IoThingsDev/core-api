package models

type Card struct {
	Id string `json:"id" bson:"id" binding:"required"`
}
