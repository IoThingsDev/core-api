package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http"
	"github.com/dernise/pushpal-api/models"
	"gopkg.in/mgo.v2/bson"
)


func TestCreateAccount(t *testing.T) {
	api := SetupRouterAndDatabase()
	defer api.Database.Session.Close()

	//Missing field
	parameters := []byte(`
	{
		"username":"dernise",
		"password":"test",
		"firstname":"maxence",
		"lastname": "henneron"
	}`)
	resp := SendRequest(api, parameters, "POST", "/v1/users/")
	assert.Equal(t, resp.Code, http.StatusBadRequest)

	//Everything is fine
	parameters = []byte(`
	{
		"username":"dernise",
		"email":"maxence.henneron@icloud.com",
		"password":"test",
		"firstname":"maxence",
		"lastname": "henneron"
	}`)
	resp = SendRequest(api, parameters, "POST", "/v1/users/")
	assert.Equal(t, resp.Code, http.StatusCreated)

	// User already exists
	resp = SendRequest(api, parameters, "POST", "/v1/users/")
	assert.Equal(t, resp.Code, http.StatusConflict)

	// Test activation
	user := models.User{}
	err := api.Database.C(models.UsersCollection).Find(bson.M{"email": "maxence.henneron@icloud.com"}).One(&user)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, user.Active, false)
	resp = SendRequest(api, nil, "GET", "/v1/users/" + user.Id.Hex() + "/activate/" + user.ActivationKey)

	//Update user informations
	err = api.Database.C(models.UsersCollection).Find(bson.M{"email": "maxence.henneron@icloud.com"}).One(&user)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, resp.Code, http.StatusOK)
	assert.Equal(t, user.Active, true)

	//Activation key isn't right
	resp = SendRequest(api, nil, "GET", "/v1/users/" + user.Id.Hex() + "/activate/fakeKey")
	assert.Equal(t, resp.Code, http.StatusNotFound)

	//Unknown user
	resp = SendRequest(api, nil, "GET", "/v1/users/" + bson.NewObjectId().Hex() + "/activate/fakeKey")
	assert.Equal(t, resp.Code, http.StatusNotFound)
}

