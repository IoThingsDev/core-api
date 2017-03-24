package server

import (
	"github.com/dernise/base-api/controllers"
	"github.com/dernise/base-api/middlewares"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "You successfully reached the base API."})
}

func (a API) SetupRouter() {
	router := a.Router

	router.Use(middlewares.ErrorMiddleware())
	userController := controllers.NewUserController(a.Database, a.EmailSender, a.Config)

	v1 := router.Group("/v1")
	{
		v1.GET("/", Index)
		users := v1.Group("/users")
		{
			users.GET("/", userController.GetUsers)
			users.POST("/", userController.CreateUser)
			users.POST("/requestResetPassword", userController.ResetPasswordRequest)
		}

		user := v1.Group("/user")
		{
			user.GET("/:id", userController.GetUser)
			user.GET("/:id/activate/:activationKey", userController.ActivateUser)
			//users.GET("/:id/reset/:resetKey", userController.FormResetPassword)
			user.POST("/:id/reset/", userController.ResetPassword)
		}

		authentication := v1.Group("/auth")
		{
			authController := controllers.NewAuthController(a.Database, a.Config)
			authentication.POST("/", authController.Authentication)
		}

		authorized := v1.Group("/authorized")
		{
			authorized.Use(middlewares.AuthMiddleware())
			billing := authorized.Group("/billing")
			{
				billingController := controllers.NewBillingController(a.Database, a.EmailSender, a.Config)
				billing.POST("/", billingController.CreateTransaction)
			}
		}
	}
}
