package main

import (
	"gopkg.in/mgo.v2"
	"github.com/dernise/pushpal-api/server"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/asaskevich/govalidator"
)

func main() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	govalidator.SetFieldsRequiredByDefault(true)

	database := session.DB("pushpal")

	api := server.API{ Router: gin.Default() }
	api.SetupRouter(database)
	api.Router.Run(":4000")
}
