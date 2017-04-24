package controllers

import (
	"net/http"

	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
	"github.com/dernise/base-api/services"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
)

type SubscriptionController struct {
	mgo    *mgo.Database
	config *viper.Viper
	redis  *services.Redis
}

func NewSubscriptionController(mgo *mgo.Database, config *viper.Viper, redis *services.Redis) SubscriptionController {
	return SubscriptionController{
		mgo,
		config,
		redis,
	}
}

func (sc SubscriptionController) CreateSubscription(c *gin.Context) {
	user := models.GetUserFromContext(c)

	plan := models.Plan{}
	if err := c.Bind(&plan); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	_, err := sub.New(&stripe.SubParams{
		Customer: user.StripeId,
		Plan:     plan.Name,
	})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("subscription_creation_failed", "Failed to subscribe the user to this plan"))
		return
	}

	c.JSON(http.StatusOK, nil)
}
