package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	api := SetupApi()
	defer api.Database.Session.Close()

	user, jwtToken := CreateUserAndGenerateToken(api)

	parameters := []byte(`
	{
		"userId":"` + user.Id.Hex() + `",
		"amount":1000
	}`)

	resp := SendRequestWithToken(api, parameters, "POST", "/v1/authorized/billing/", jwtToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}
