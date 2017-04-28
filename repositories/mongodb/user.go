package mongodb

import (
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func (db *mongo) CreateUser(user *models.User) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user.Email = strings.ToLower(user.Email)
	if count, _ := users.Find(bson.M{"email": user.Email}).Count(); count > 0 {
		return helpers.ErrorWithCode("user_already_exists", "User already exists")
	}

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return err
	}

	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	if err != nil {
		return helpers.ErrorWithCode("password_encryption_failed", "Failed to generate the encrypted password")
	}

	user.Active = false
	user.ActivationKey = helpers.RandomString(20)
	user.StripeId = ""
	user.Id = bson.NewObjectId()
	user.Admin = false

	err = users.Insert(user)
	if err != nil {
		return helpers.ErrorWithCode("user_creation_failed", "Failed to insert the user in the mongobase")
	}

	return nil
}

func (db *mongo) GetUser(id string) (*models.User, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user := &models.User{}
	err := users.FindId(bson.ObjectIdHex(id)).One(user)

	return user, err
}

func (db *mongo) ActivateUser(activationKey string, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	err := users.Update(bson.M{"$and": []bson.M{{"_id": bson.ObjectIdHex(id)}, {"activationKey": activationKey}}}, bson.M{"$set": bson.M{"active": true}})
	if err != nil {
		return helpers.ErrorWithCode("user_activation_failed", "Couldn't find the user to activate")
	}
	return nil
}