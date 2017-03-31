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

func (a *API) SetupRouter() {
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
			users.POST("/", userController.CreateUser)
			users.GET("/:id/activate/:activationKey", userController.ActivateUser)
			users.POST("/:id/reset_password", userController.ResetPassword)

			users.GET("/:id", userController.GetUser).Use(middlewares.AuthMiddleware())
		}

		cards := v1.Group("/cards")
		{
			cards.Use(middlewares.AuthMiddleware())
			cardController := controllers.NewCardController(a.Database, a.Config)
			cards.POST("/", cardController.AddCard)
			cards.GET("/", cardController.GetCards)
			cards.PUT("/:id/set_default", cardController.SetDefaultCard)
			cards.DELETE("/:id", cardController.DeleteCard)
		}

		authentication := v1.Group("/auth")
		{
			authController := controllers.NewAuthController(a.Database, a.Config)
			authentication.POST("/", authController.Authentication)
		}

		billing := v1.Group("/billing")
		{
			billing.Use(middlewares.AuthMiddleware())
			billingController := controllers.NewBillingController(a.Database, a.EmailSender, a.Config)
			billing.POST("/", billingController.CreateTransaction)
		}
	}
}
