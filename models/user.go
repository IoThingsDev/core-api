package models

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id        bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	Firstname string   `json:"firstname" bson:"firstname"`
	Lastname  string   `json:"lastname" bson:"lastname"`
	Password  string   `json:"password" bson:"password"`
	Email     string `json:"email" bson:"email" valid:"email"`
	Username  string `json:"username" bson:"username"`
}

const UsersCollection = "users"
