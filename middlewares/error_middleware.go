package middlewares

import (
	"github.com/dernise/base-api/helpers"
	"gopkg.in/gin-gonic/gin.v1"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errorToPrint := c.Errors.Last()
		if errorToPrint != nil {
			original, ok := errorToPrint.Err.(helpers.Error)
			if ok {
				c.JSON(-1, gin.H{
					"status":  "error",
					"message": original.Message,
					"code":    original.Code,
				})
			} else {
				c.JSON(-1, gin.H{
					"status":  "error",
					"message": errorToPrint.Error(),
					"code":    "unknown",
				})
			}
		}
	}
}
