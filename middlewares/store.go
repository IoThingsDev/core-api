package middlewares

import (
	"github.com/dernise/base-api/store"
	"github.com/dernise/base-api/store/mongodb"
	"gopkg.in/gin-gonic/gin.v1"
	mgo "gopkg.in/mgo.v2"
)

func StoreMiddleware(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		store.ToContext(c, mongodb.New(db))
		c.Next()
	}
}
