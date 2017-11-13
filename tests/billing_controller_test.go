package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*func TestCreateTransaction(t *testing.T) {
	parameters := []byte(`
	{
		"userId":"` + user.Id.Hex() + `",
		"amount":100
	}`)

	resp := SendRequestWithToken(parameters, "POST", "/v1/billing/", authToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}*/

func TestCreatePlan(t *testing.T) {
	parameters := []byte(`
	{
		"id":"best-plan",
		"amount": 1000,
		"name": "The best plan for you!",
		"interval": "month",
		"currency": "eur",
		"metadata": {"description":"plan that allows you to use one app and one user"}
	}`)

	resp := SendRequestWithToken(parameters, "POST", "/v1/billing/plans/", authToken)
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestGetPlans(t *testing.T) {
	resp := SendRequestWithToken(nil, "GET", "/v1/billing/plans/", authToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}
