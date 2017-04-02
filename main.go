package main

import (
	"github.com/dernise/base-api/server"
	"github.com/dernise/base-api/services"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	api := &server.API{Router: gin.Default(), Config: viper.New()}

	// Configuration setup
	err := api.SetupViper()
	if err != nil {
		panic(err)
	}

	// Email sender setup
	api.EmailSender = services.NewSendGridEmailSender(api.Config)

	// Database setup
	session, err := api.SetupDatabase()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = api.SetupIndexes()
	if err != nil {
		panic(err)
	}

	// Stripe setup
	services.SetStripeKeyAndBackend(api.Config)

	// Redis setup
	api.SetupRedis()

	// Router setup
	api.SetupRouter()
	api.Router.Run(api.Config.GetString("host_address"))
}
