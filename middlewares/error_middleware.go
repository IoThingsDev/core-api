package middlewares

import "gopkg.in/gin-gonic/gin.v1"

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errorToPrint := c.Errors.Last()
		if errorToPrint != nil {
			c.JSON(-1, gin.H {
				"status": "error",
				"message": errorToPrint.Error(),
			})
		}
	}
}