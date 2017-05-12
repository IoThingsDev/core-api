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

	router.Use(middlewares.StoreMiddleware(a.Database))
	router.Use(middlewares.ConfigMiddleware(a.Config))
	router.Use(middlewares.RedisMiddleware(a.Redis))

	router.Use(middlewares.EmailMiddleware(a.EmailSender))
	router.Use(middlewares.RateMiddleware())

	authMiddleware := middlewares.AuthMiddleware()

	v1 := router.Group("/v1")
	{
		v1.GET("/", Index)
		userController := controllers.NewUserController()
		//v1.POST("/reset_password", userController.ResetPasswordRequest)
		users := v1.Group("/users")
		{
			users.POST("/", userController.CreateUser)
			users.GET("/:id/activate/:activationKey", userController.ActivateUser)
			//users.POST("/:id/reset_password", userController.ResetPassword)

			users.Use(authMiddleware)
			users.GET("/:id", userController.GetUser)
		}

		cards := v1.Group("/cards")
		{
			cards.Use(authMiddleware)
			cardController := controllers.NewCardController()
			cards.POST("/", cardController.AddCard)
			cards.GET("/", cardController.GetCards)
			cards.PUT("/:id/set_default", cardController.SetDefaultCard)
			cards.DELETE("/:id", cardController.DeleteCard)
		}

		authentication := v1.Group("/auth")
		{
			authController := controllers.NewAuthController()
			authentication.POST("/", authController.Authentication)
		}

		billing := v1.Group("/billing")
		{
			billingController := controllers.NewBillingController()
			billing.Use(authMiddleware)

			plans := billing.Group("/plans")
			{
				plans.GET("/", billingController.GetPlans)
				plans.POST("/", middlewares.AdminMiddleware(), billingController.CreatePlan)
			}

			subscriptionController := controllers.NewStripeSubscriptionController()
			subscriptions := billing.Group("/subscriptions")
			{
				subscriptions.POST("/", subscriptionController.CreateSubscription)
				subscriptions.GET("/", subscriptionController.GetSubscriptions)
				subscriptions.DELETE("/:id", subscriptionController.DeleteSubscription)
			}
		}
	}
}
