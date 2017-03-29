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
	Active        bool          `json:"active" bson:"active" valid:"-"`
	ActivationKey string        `json:"activationKey" bson:"activationKey" valid:"-"`
	ResetKey      string        `json:"resetKey" bson:"resetKey" valid:"-"`
	StripeId      string        `json:"stripeId" bson:"stripeId" valid:"-"`
	Cards         []Card        `json:"cards" bson:"cards" valid:"-"`
}

type SanitizedUser struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Firstname string        `json:"firstname" bson:"firstname"`
	Lastname  string        `json:"lastname" bson:"lastname"`
	Email     string        `json:"email" bson:"email"`
}

func (uc User) Sanitize() SanitizedUser {
	return SanitizedUser{uc.Id, uc.Firstname, uc.Lastname, uc.Email}
}

const UsersCollection = "users"
