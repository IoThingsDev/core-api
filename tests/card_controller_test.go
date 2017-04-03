package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestAddCard(t *testing.T) {
	api := SetupApi()
	defer api.Database.Session.Close()

	_, jwtToken := CreateUserAndGenerateToken(api)

	parameters := []byte(`
	{
		"token":"TestToken"
	}`)

	resp := SendRequestWithToken(api, parameters, "POST", "/v1/cards/", jwtToken)
	assert.Equal(t, http.StatusCreated, resp.Code)

	parameters = []byte(`
	{
		"oken":"TestToken"
	}`)

	resp = SendRequestWithToken(api, parameters, "POST", "/v1/cards/", jwtToken)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetCards(t *testing.T) {
	api := SetupApi()
	defer api.Database.Session.Close()

	_, jwtToken := CreateUserAndGenerateToken(api)

	resp := SendRequestWithToken(api, nil, "GET", "/v1/cards/", jwtToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDefaultCard(t *testing.T) {
	api := SetupApi()
	defer api.Database.Session.Close()

	_, jwtToken := CreateUserAndGenerateToken(api)

	resp := SendRequestWithToken(api, nil, "PUT", "/v1/cards/testId/set_default", jwtToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDeleteCard(t *testing.T) {
	api := SetupApi()
	defer api.Database.Session.Close()

	_, jwtToken := CreateUserAndGenerateToken(api)

	resp := SendRequestWithToken(api, nil, "DELETE", "/v1/cards/"+bson.NewObjectId().Hex(), jwtToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}
