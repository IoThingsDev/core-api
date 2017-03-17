package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	api := SetupRouterAndDatabase()

	//Missing field
	parameters := []byte(`
	{
		"username":"dernise",
		"password":"test",
		"firstname":"maxence",
		"lastname": "henneron"
	}`)
	resp := SendRequest(api, parameters, "POST", "/v1/users/")
	assert.Equal(t, resp.Code, 400)

	//Everything is fine
	parameters = []byte(`
	{
		"username":"dernise",
		"email":"maxenjce.henneron@icloud.com",
		"password":"test",
		"firstname":"maxence",
		"lastname": "henneron"
	}`)
	resp = SendRequest(api, parameters, "POST", "/v1/users/")
	assert.Equal(t, resp.Code, 201)

	// User already exists
	resp = SendRequest(api, parameters, "POST", "/v1/users/")
	assert.Equal(t, resp.Code, 409)
}

