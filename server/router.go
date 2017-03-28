package server

import (
	"net/http"
	"time"

	"github.com/dernise/base-api/controllers"
	"github.com/dernise/base-api/middlewares"
	"gopkg.in/gin-gonic/gin.v1"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "You successfully reached the base API."})
}

func (a API) SetupRouter() {
	router := a.Router

	router.Use(middlewares.ErrorMiddleware())

	router.Use(middlewares.CorsMiddleware(middlewares.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	v1 := router.Group("/v1")
	{
		v1.GET("/", Index)
		userController := controllers.NewUserController(a.Database, a.EmailSender, a.Config)
		v1.POST("/reset_password", userController.ResetPasswordRequest)
		users := v1.Group("/users")
		{
			users.GET("/", userController.GetUsers)
			users.POST("/", userController.CreateUser)
			users.GET("/:id", userController.GetUser)
			users.GET("/:id/activate/:activationKey", userController.ActivateUser)
			//users.GET("/:id/reset/:resetKey", userController.FormResetPassword)
			users.POST("/:id/reset_password", userController.ResetPassword)
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
