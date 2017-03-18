package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http"
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
}

