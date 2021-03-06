package server

import (
	"net/http"
	"time"

	"github.com/IoThingsDev/api/controllers"
	"github.com/IoThingsDev/api/middlewares"

	"gopkg.in/gin-gonic/gin.v1"
)

func Index(c *gin.Context) {
	//TODO: Implement some Viper config customizing based on BASEAPI_SENDGRID_NAME
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "You successfully reached ThingsApi."})
}

func (a *API) SetupRouter() {
	router := a.Router

	router.LoadHTMLGlob("templates/pages")

	router.Use(middlewares.ErrorMiddleware())

	router.Use(middlewares.CorsMiddleware(middlewares.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, X-Requested-With, Accept, Token",
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
		//TODO : Implement robots.txt :
		//User-Agent: *
		//Disallow: /
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

		/*		cards := v1.Group("/cards")
				{
					cards.Use(authMiddleware)
					cardController := controllers.NewCardController()
					cards.POST("/", cardController.AddCard)
					cards.GET("/", cardController.GetCards)
					cards.PUT("/:id/set_default", cardController.SetDefaultCard)
					cards.DELETE("/:id", cardController.DeleteCard)
				}
		*/
		sigfox := v1.Group("/sigfox")
		{
			sigfoxController := controllers.NewSigfoxController()
			sigfox.POST("/messages", sigfoxController.CreateMessage)

			sigfox.POST("/messages/import", sigfoxController.ImportMessage)
			sigfox.POST("/locations/import", sigfoxController.ImportLocation)


			locationController := controllers.NewLocationController()
			sigfox.POST("/locations", locationController.CreateLocation)

			sigfox.POST("/atlas", locationController.CreateLocation)
		}

		devices := v1.Group("/devices")
		{
			devices.Use(authMiddleware)
			deviceController := controllers.NewDeviceController()
			devices.GET("/", deviceController.GetDevices)
			devices.POST("/", deviceController.CreateDevice)
			devices.PUT("/:id", deviceController.UpdateDevice)
			devices.GET("/:id", deviceController.GetDevice)
			devices.DELETE("/:id", deviceController.DeleteDevice)
			devices.GET("/:id/locations", deviceController.GetAllDeviceLocations)
			devices.GET("/:id/messages", deviceController.GetAllDeviceMessages)
			devices.GET("/:id/lastLocations", deviceController.GetLastDeviceLocations)
			devices.GET("/:id/lastMessages", deviceController.GetLastDeviceMessages)
		}

		messages := v1.Group("/messages")
		{
			messages.Use(authMiddleware)
			messagesController := controllers.NewSigfoxController()
			messages.GET("/", messagesController.GetLastDevicesSigfoxMessages)
		}

		locations := v1.Group("/locations")
		{
			locations.Use(authMiddleware)
			locationController := controllers.NewLocationController()
			locations.GET("/", locationController.GetLastDevicesLocations)

		}

		authentication := v1.Group("/auth")
		{
			authController := controllers.NewAuthController()
			authentication.POST("/", authController.Authentication)
			authentication.OPTIONS("/", authController.Preflight)
			authentication.Use(authMiddleware)
			authentication.GET("/logout", authController.LogOut)
		}

		/*billing := v1.Group("/billing")
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
		}*/
	}
}
