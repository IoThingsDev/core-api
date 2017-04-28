package middlewares

import (
	"github.com/dernise/base-api/repositories"
	"github.com/dernise/base-api/repositories/mongodb"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
)

func StoreMiddleware(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		repositories.ToContext(c, mongodb.New(db))
		c.Next()
	}
}
