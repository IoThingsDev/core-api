package models

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id            bson.ObjectId `json:"id" bson:"_id,omitempty" valid:"-"`
	Firstname     string        `json:"firstname" bson:"firstname" valid:"required"`
	Lastname      string        `json:"lastname" bson:"lastname" valid:"required"`
	Password      string        `json:"password" bson:"password" valid:"required"`
	Email         string        `json:"email" bson:"email" valid:"email,required"`
	Username      string        `json:"username" bson:"username" valid:"required"`
	Active        bool          `json:"active" bson:"active" valid:"-"`
	ActivationKey string        `json:"activationKey" bson:"activationKey" valid:"-"`
	StripeId      string        `json:"stripeId" bson:"stripeId" valid:"-"`
}

const UsersCollection = "users"
