package main

import (
	"testing"
	"fmt"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"net/http"
)



func TestHomePage(t *testing.T) {
	api := SetupRouterAndDatabase()

	req, err := http.NewRequest("GET", "/v1/", nil)
	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}