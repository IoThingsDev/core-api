package main

import (
	"github.com/dernise/base-api/server"
	"github.com/dernise/base-api/services"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	api := &server.API{Router: gin.Default(), Config: viper.New()}

	err := api.SetupViper()
	if err != nil {
		panic(err)
	}

	api.EmailSender = services.NewSendGridEmailSender(api.Config)

	session, err := api.SetupDatabase()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = api.SetupIndexes()
	if err != nil {
		panic(err)
	}

	services.SetStripeKeyAndBackend(api.Config)

	api.SetupRouter()
	api.Router.Run(api.Config.GetString("host_address"))
}
