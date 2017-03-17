package tests

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/dernise/pushpal-api/server"
	"github.com/spf13/viper"
	"bytes"
	"net/http/httptest"
	"net/http"
)

func SetupRouterAndDatabase() *server.API {
	api := server.API{ Router: gin.Default(), Config: viper.New() }
	api.SetupViper("test")
	session, err := mgo.Dial(api.Config.GetString("database.address"))
	if err != nil {
		panic(err)
	}

	api.Database = session.DB(api.Config.GetString("database.dbName"))
	api.Database.DropDatabase()
	api.SetupRouter()

	return &api
}

func SendRequest(api *server.API, parameters []byte, method string, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(parameters))
	req.Header.Add("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	return resp
}