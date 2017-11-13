package models

type LoginToken struct {
	Id         string `json:"id" bson:"_id"`
	Ip         string `json:"ip" bson:"ip"`
	CreatedAt  int64  `json:"created_at" bson:"created_at"`
	LastAccess int64  `json:"last_access" bson:"last_access"`
}
