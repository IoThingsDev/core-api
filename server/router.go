package server

import (
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"github.com/dernise/pushpal-api/controllers"
	"github.com/dernise/pushpal-api/middlewares"
	"net/http"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "You successfully reached the Pushpal API."})
}

func (a API) SetupRouter(database *mgo.Database) {
	router := a.Router

	router.Use(middlewares.ErrorMiddleware())

	v1 := router.Group("/v1")
	{
		v1.GET("/", Index)
		users := v1.Group("/users")
		{
			userController := controllers.NewUserController(database)
			users.GET("/", userController.GetUsers)
			users.GET("/:id", userController.GetUser)
			users.POST("/", userController.CreateUser)
		}

		authentication := v1.Group("/auth")
		{
			authController := controllers.NewAuthController(database)
			authentication.POST("/", authController.Authentication)
		}

		authorized := v1.Group("/authorized").Use(middlewares.AuthMiddleware())
		{
			authorized.GET("/", Index)
		}
	}
}

