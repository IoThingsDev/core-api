package models

import (
	"net/http"
	"strings"

	"github.com/adrien3d/things-api/helpers"
	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id            string       `json:"id" bson:"_id,omitempty" valid:"-"`
	Firstname     string       `json:"firstname" bson:"firstname"`
	Lastname      string       `json:"lastname" bson:"lastname"`
	Password      string       `json:"password" bson:"password" valid:"required"`
	Email         string       `json:"email" bson:"email" valid:"email,required"`
	Phone         string       `json:"phone" bson:"phone"`
	Active        bool         `json:"active" bson:"active"`
	ActivationKey string       `json:"activationKey" bson:"activationKey"`
	ResetKey      string       `json:"resetKey" bson:"resetKey"`
	StripeId      string       `json:"stripeId" bson:"stripeId"`
	Admin         bool         `json:"admin" bson:"admin"`
	Tokens        []LoginToken `json:"tokens" bson:"tokens"`
}

type SanitizedUser struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	Firstname string `json:"firstname" bson:"firstname"`
	Lastname  string `json:"lastname" bson:"lastname"`
	Email     string `json:"email" bson:"email"`
}

func (user *User) Sanitize() SanitizedUser {
	return SanitizedUser{user.Id, user.Firstname, user.Lastname, user.Email}
}

func (user *User) BeforeCreate() error {
	user.Active = false
	user.ActivationKey = helpers.RandomString(20)
	user.StripeId = ""
	user.Admin = false
	user.Email = strings.ToLower(user.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "encryption_failed", "Failed to generate the crypted password")
	}
	user.Password = string(hashedPassword)

	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		return helpers.NewError(http.StatusBadRequest, "input_not_valid", err.Error())
	}

	return nil
}

func (user *User) HasToken(tokenId string) (int, bool) {
	for index, token := range user.Tokens {
		if token.Id == tokenId {
			return index, true
		}
	}
	return -1, false
}

const UsersCollection = "users"
