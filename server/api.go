package server

import (
	"github.com/dernise/base-api/services"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
)

type API struct {
	Router      *gin.Engine
	Config      *viper.Viper
	Database    *mgo.Database
	EmailSender services.EmailSender
}
