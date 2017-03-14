package models

type Address struct {
	City  string `json:"city" bson:"city"`
	State string `json:"state" bson:"state"`
}
