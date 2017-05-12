package mongodb

import (
	"net/http"

	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/helpers/params"
	"github.com/dernise/base-api/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *mongo) CreateUser(user *models.User) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user.Id = bson.NewObjectId().Hex()
	err := user.BeforeCreate()
	if err != nil {
		return err
	}

	if count, _ := users.Find(bson.M{"email": user.Email}).Count(); count > 0 {
		return helpers.NewError(http.StatusConflict, "user_already_exists", "User already exists")
	}

	err = users.Insert(user)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_creation_failed", "Failed to insert the user in the database")
	}

	return nil
}

func (db *mongo) FindUserById(id string) (*models.User, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user := &models.User{}
	err := users.FindId(id).One(user)

	return user, err
}

func (db *mongo) FindUser(params params.M) (*models.User, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user := &models.User{}

	err := users.Find(params).One(user)

	return user, err
}

func (db *mongo) ActivateUser(activationKey string, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	err := users.Update(bson.M{"$and": []bson.M{{"_id": id}, {"activationKey": activationKey}}}, bson.M{"$set": bson.M{"active": true}})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_activation_failed", "Couldn't find the user to activate")
	}
	return nil
}

func (db *mongo) UpdateUser(user *models.User, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	err := users.UpdateId(user.Id, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_update_failed", "Failed to update the user")
	}

	return nil
}
