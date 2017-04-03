package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateMiddleware(t *testing.T) {
	api := SetupApi()
	defer api.Database.Session.Close()

	req, err := http.NewRequest("GET", "/v1/", nil)
	if err != nil {
		fmt.Println(err)
	}

	api.Config.Set("redis_should_activate_rates", true)

	for i := 0; i < 5; i++ {
		resp := httptest.NewRecorder()
		api.Router.ServeHTTP(resp, req)
		assert.Equal(t, resp.Code, 200)
	}

	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusTooManyRequests)

	api.Config.Set("redis_should_activate_rates", false)
}
