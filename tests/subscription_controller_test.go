package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSubscription(t *testing.T) {
	parameters := []byte(`
	{
		"id":"fakeplanid"
	}`)

	resp := SendRequestWithToken(parameters, "POST", "/v1/billing/subscriptions/", authToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetSubscriptions(t *testing.T) {
	resp := SendRequestWithToken(nil, "GET", "/v1/billing/subscriptions/", authToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDeleteSubscription(t *testing.T) {
	resp := SendRequestWithToken(nil, "DELETE", "/v1/billing/subscriptions/fakeid", authToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}
