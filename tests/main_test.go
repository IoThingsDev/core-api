package tests

import (
	"os"
	"testing"

	"github.com/dernise/base-api/models"
	"github.com/dernise/base-api/server"
)

var api *server.API
var user *models.User
var jwtToken string

func TestMain(m *testing.M) {
	api = SetupApi()
	user, jwtToken = CreateUserAndGenerateToken()
	retCode := m.Run()
	api.Database.Session.Close()
	os.Exit(retCode)
}
