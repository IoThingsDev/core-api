package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	parameters := []byte(`
	{
		"userId":"` + user.Id.Hex() + `",
		"amount":100
	}`)

	resp := SendRequestWithToken(parameters, "POST", "/v1/billing/", jwtToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}
