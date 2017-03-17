package server

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/spf13/viper"
)

type API struct {
	Router *gin.Engine
	Config *viper.Viper
}
