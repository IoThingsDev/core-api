package controllers

import (
	"net/http"

	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
	"github.com/dernise/base-api/services"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
	"gopkg.in/gin-gonic/gin.v1"
)

type StripeSubscriptionController struct {
	redis *services.Redis
}

func NewStripeSubscriptionController(redis *services.Redis) StripeSubscriptionController {
	return StripeSubscriptionController{
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
		Plan:     plan.Id,
	})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("subscription_creation_failed", err.Error()))
		return
	}

	sc.redis.InvalidateObject(user.StripeId + "-subscriptions")

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

	if _, err := sub.Cancel(subscriptionId, &stripe.SubParams{
		Customer: user.StripeId,
	}); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("delete_subscription_failed", "Failed to delete the subscription"))
		return
	}

	sc.redis.InvalidateObject(user.StripeId + "-subscriptions")

	c.JSON(http.StatusOK, nil)
}

func (sc StripeSubscriptionController) getSubscriptions(stripeId string) []*stripe.Sub {
	subscriptions := []*stripe.Sub{}
	err := sc.redis.GetValueForKey(stripeId+"-subscriptions", subscriptions)
	if err != nil {
		params := &stripe.SubListParams{}
		params.Customer = stripeId
		i := sub.List(params)
		for i.Next() {
			subscriptions = append(subscriptions, i.Sub())
		}

		sc.redis.SetValueForKey(stripeId+"-subscriptions", subscriptions)
	}

	return subscriptions
}
