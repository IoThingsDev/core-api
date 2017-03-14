package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id        bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	Firstname string   `json:"firstname" bson:"firstname"`
	Lastname  string   `json:"lastname" bson:"lastname"`
	Address   *Address `json:"address" bson:"address"`
}

const UsersCollection = "users"