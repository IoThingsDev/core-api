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

type StripeSubscriptionController struct {
	mgo    *mgo.Database
	config *viper.Viper
	redis  *services.Redis
}

func NewStripeSubscriptionController(mgo *mgo.Database, config *viper.Viper, redis *services.Redis) StripeSubscriptionController {
	return StripeSubscriptionController{
		mgo,
		config,
		redis,
	}
}

func (sc StripeSubscriptionController) CreateSubscription(c *gin.Context) {
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

func (sc StripeSubscriptionController) GetSubscriptions(c *gin.Context) {
	user := models.GetUserFromContext(c)

	subscriptions := sc.getSubscriptions(user.StripeId)

	c.JSON(http.StatusOK, gin.H{"subscriptions": subscriptions})
}

func (sc StripeSubscriptionController) DeleteSubscription(c *gin.Context) {
	user := models.GetUserFromContext(c)

	subscriptionId := c.Param("id")

	subscriptions := sc.getSubscriptions(user.StripeId)

	found := false
	for _, s := range subscriptions {
		if subscriptionId == s.ID {
			found = true
			break
		}
	}

	if found == false {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "The subscription id is wrong"))
		return
	}

	if _, err := sub.Cancel(subscriptionId, nil); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("delete_subscription_failed", "Failed to delete the subscription"))
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (sc StripeSubscriptionController) getSubscriptions(stripeId string) []*stripe.Sub {
	subscriptions := []*stripe.Sub{}
	params := &stripe.SubListParams{}
	params.Customer = stripeId
	i := sub.List(params)
	for i.Next() {
		subscriptions = append(subscriptions, i.Sub())
	}

	return subscriptions
}
