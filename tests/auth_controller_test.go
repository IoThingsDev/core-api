package tests

import (
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"bytes"
	"log"
)

func TestCreateAccount(t *testing.T) {
	api := SetupRouterAndDatabase()

	parameters := []byte(`{"username":"dernise", "email":"maxence.henneron@icloud.com", "password":"test", "firstname":"maxence", "lastname": "henneron"}`)

	req, err := http.NewRequest("POST", "/v1/users/", bytes.NewBuffer(parameters))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
	}

	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 201)
}

