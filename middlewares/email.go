package middlewares

import (
	"github.com/dernise/base-api/services"
	"gopkg.in/gin-gonic/gin.v1"
)

func EmailMiddleware(emailSender services.EmailSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("emailSender", emailSender)
		c.Next()
	}
}
