	package middlewares

import (
	"github.com/IoThingsDev/api/services"

	"gopkg.in/gin-gonic/gin.v1"
)

func RedisMiddleware(redis *services.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redis", redis)
		c.Next()
	}
}
