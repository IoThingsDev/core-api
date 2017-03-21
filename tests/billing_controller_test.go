package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func testCreateTransaction(t *testing.T) {
	api := SetupRouterAndDatabase()
	defer api.Database.Session.Close()

	TestCreateAccount(t)

	parameters := []byte(`
	{
		"userId":"58d12bfe373d36799a0cf187",
		"amount":1000,
		"cardToken":"tok_1A05hsAPpGNju8w2s6LHtFYG"
	}`)

	resp := SendRequest(api, parameters, "POST", "/v1/authorized/billing")
	assert.Equal(t, resp.Code, http.StatusBadRequest)
}
