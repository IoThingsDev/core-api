package server

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"github.com/dernise/base-api/services"
)

type API struct {
	Router *gin.Engine
	Config *viper.Viper
	Database *mgo.Database
	EmailSender services.EmailSender
}
