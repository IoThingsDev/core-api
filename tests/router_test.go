package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomePage(t *testing.T) {
	api := SetupRouterAndDatabase()
	defer api.Database.Session.Close()

	req, err := http.NewRequest("GET", "/v1/", nil)
	if err != nil {
		fmt.Println(err)
	}

	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}
