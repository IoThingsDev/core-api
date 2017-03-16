package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/dernise/pushpal-api/server"
)

func SetupRouterAndDatabase() *server.API {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	database := session.DB("pushpal-tests")

	api := server.API{Router: gin.Default() }
	api.SetupRouter(database)

	return &api
}
