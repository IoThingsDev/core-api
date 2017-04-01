package tests

import (
	"net/http"
	"testing"

	"github.com/dernise/base-api/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestCreateAccount(t *testing.T) {
	api := SetupApi()
	defer api.Database.Session.Close()

	//Missing field
	parameters := []byte(`
	{
		"password":"test",
		"firstname":"maxence",
		"lastname": "henneron"
	}`)
	resp := SendRequest(api, parameters, "POST", "/v1/users/")
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	//Everything is fine
	parameters = []byte(`
	{
		"email":"maxence.henneron@icloud.com",
		"password":"test",
		"firstname":"maxence",
		"lastname": "henneron"
	}`)
	resp = SendRequest(api, parameters, "POST", "/v1/users/")
	assert.Equal(t, http.StatusCreated, resp.Code)

	// User already exists
	resp = SendRequest(api, parameters, "POST", "/v1/users/")
	assert.Equal(t, http.StatusConflict, resp.Code)

	// Duplicate email
	parameters = []byte(`
	{
		"email":"mAxEnce.henneron@icloud.com",
		"password":"test",
		"firstname":"maxence",
		"lastname": "henneron"
	}`)
	resp = SendRequest(api, parameters, "POST", "/v1/users/")
	assert.Equal(t, http.StatusConflict, resp.Code)

	// Test activation
	user := models.User{}
	err := api.Database.C(models.UsersCollection).Find(bson.M{"email": "maxence.henneron@icloud.com"}).One(&user)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, user.Active, false)
	resp = SendRequest(api, nil, "GET", "/v1/users/"+user.Id.Hex()+"/activate/"+user.ActivationKey)

	//Update user information
	err = api.Database.C(models.UsersCollection).Find(bson.M{"email": "maxence.henneron@icloud.com"}).One(&user)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, user.Active, true)

	//Activation key isn't right
	resp = SendRequest(api, nil, "GET", "/v1/users/"+user.Id.Hex()+"/activate/fakeKey")
	assert.Equal(t, http.StatusNotFound, resp.Code)

	//Unknown user
	resp = SendRequest(api, nil, "GET", "/v1/users/"+bson.NewObjectId().Hex()+"/activate/fakeKey")
	assert.Equal(t, http.StatusNotFound, resp.Code)

	//Delete user
	resp = SendRequest(api, nil, "DELETE", "/v1/users/"+user.Id.Hex())
	assert.Equal(t, http.StatusOK, resp.Code)
}
