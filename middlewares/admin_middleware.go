package middlewares

import (
	"net/http"

	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
	"gopkg.in/gin-gonic/gin.v1"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface := c.MustGet("currentUser")
		user := userInterface.(models.User)

		if !user.Admin {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("admin_required", "The user is not administrator"))
			return
		}
	}
}
