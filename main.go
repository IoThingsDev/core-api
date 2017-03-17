package main

import (
	"gopkg.in/mgo.v2"
	"github.com/dernise/pushpal-api/server"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
)

func main() {
	api := server.API{ Router: gin.Default(), Config: viper.New() }
	api.SetupViper()

	session, err := mgo.Dial(api.Config.GetString("database.address"))
	if err != nil {
		panic(err)
	}
	defer session.Close()
	database := session.DB(api.Config.GetString("database.dbName"))

	govalidator.SetFieldsRequiredByDefault(true)

	api.SetupRouter(database)
	api.Router.Run(api.Config.GetString("host.port"))
}
