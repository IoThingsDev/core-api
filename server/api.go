package server

import (
	"github.com/IoThingsDev/api/services"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
)

// API Structure that gathers all services
type API struct {
	Router      *gin.Engine
	Config      *viper.Viper
	Database    *mgo.Database
	EmailSender services.EmailSender
	Redis       *services.Redis
}
