package middlewares

import (
	"github.com/adrien3d/things-api/services"

	"gopkg.in/gin-gonic/gin.v1"
)

func RedisMiddleware(redis *services.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redis", redis)
		c.Next()
	}
}
