package models

import (
	"gopkg.in/gin-gonic/gin.v1"
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
	Admin         bool          `json:"admin" bson:"admin"`
}

type SanitizedUser struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Firstname string        `json:"firstname" bson:"firstname"`
	Lastname  string        `json:"lastname" bson:"lastname"`
	Email     string        `json:"email" bson:"email"`
}

func GetUserFromContext(c *gin.Context) User {
	return c.MustGet("currentUser").(User)
}

func (u *User) Sanitize() SanitizedUser {
	return SanitizedUser{u.Id, u.Firstname, u.Lastname, u.Email}
}

const UsersCollection = "users"
