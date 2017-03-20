package server

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/dernise/base-api/controllers"
	"github.com/dernise/base-api/middlewares"
	"net/http"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "You successfully reached the base API."})
}

func (a API) SetupRouter() {
	router := a.Router

	router.Use(middlewares.ErrorMiddleware())

	v1 := router.Group("/v1")
	{
		v1.GET("/", Index)
		users := v1.Group("/users")
		{
			userController := controllers.NewUserController(a.Database, a.EmailSender)
			users.GET("/", userController.GetUsers)
			users.GET("/:id", userController.GetUser)
			users.POST("/", userController.CreateUser)
			users.GET("/:id/activate/:key", userController.ActivateUser)
		}

		authentication := v1.Group("/auth")
		{
			authController := controllers.NewAuthController(a.Database)
			authentication.POST("/", authController.Authentication)
		}

		authorized := v1.Group("/authorized").Use(middlewares.AuthMiddleware())
		{
			authorized.GET("/", Index)
		}
	}
}

