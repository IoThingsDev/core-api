package middlewares

import (
	"github.com/dernise/base-api/configuration"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
)

func ConfigMiddleware(viper *viper.Viper) gin.HandlerFunc {
	return func(c *gin.Context) {
		configuration.ToContext(c, configuration.New(viper))
		c.Next()
	}
}
