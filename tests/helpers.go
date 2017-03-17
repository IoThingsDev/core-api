package tests

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/dernise/pushpal-api/server"
	"github.com/spf13/viper"
)

func SetupRouterAndDatabase() *server.API {
	api := server.API{ Router: gin.Default(), Config: viper.New() }
	api.SetupViper("test")
	session, err := mgo.Dial(api.Config.GetString("database.address"))
	if err != nil {
		panic(err)
	}
	defer session.Close()

	api.Database = session.DB(api.Config.GetString("database.dbName"))
	api.SetupRouter()
	return &api
}
