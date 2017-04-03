package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestAddCard(t *testing.T) {
	parameters := []byte(`
	{
		"token":"TestToken"
	}`)

	resp := SendRequestWithToken(parameters, "POST", "/v1/cards/", jwtToken)
	assert.Equal(t, http.StatusCreated, resp.Code)

	parameters = []byte(`
	{
		"oken":"TestToken"
	}`)

	resp = SendRequestWithToken(parameters, "POST", "/v1/cards/", jwtToken)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetCards(t *testing.T) {
	resp := SendRequestWithToken(nil, "GET", "/v1/cards/", jwtToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDefaultCard(t *testing.T) {
	resp := SendRequestWithToken(nil, "PUT", "/v1/cards/testId/set_default", jwtToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDeleteCard(t *testing.T) {
	resp := SendRequestWithToken(nil, "DELETE", "/v1/cards/"+bson.NewObjectId().Hex(), jwtToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}
