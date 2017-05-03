package controllers

import (
	"net/http"

	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/helpers/params"
	"github.com/dernise/base-api/models"
	"github.com/dernise/base-api/services"
	"github.com/dernise/base-api/store"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/customer"
	"gopkg.in/gin-gonic/gin.v1"
)

type CardController struct {
}

type Card struct {
	Token string `json:"token" binding:"required"`
}

func NewCardController() CardController {
	return CardController{}
}

func (cc CardController) AddCard(c *gin.Context) {
	user := store.Current(c)

	stripeCard := Card{}
	err := c.Bind(&stripeCard)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	if user.StripeId == "" {
		user.StripeId, err = cc.createCustomer(c, user)
		services.GetRedis(c).InvalidateObject(user.Id.Hex())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("server_error", "Failed to create the customer in our billing platform"))
			return
		}
	}

	response, err := card.New(&stripe.CardParams{
		Customer: user.StripeId,
		Token:    stripeCard.Token,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("add_card_failed", err.Error()))
		return
	}
	services.GetRedis(c).InvalidateObject(user.StripeId)

	c.JSON(http.StatusCreated, gin.H{"cards": response})
}

func (cc CardController) GetCards(c *gin.Context) {
	user := store.Current(c)

	if user.StripeId == "" {
		c.JSON(http.StatusOK, gin.H{"cards": []models.Card{}})
		return
	}

	stripeCustomer := &stripe.Customer{}

	err := services.GetRedis(c).GetValueForKey(user.StripeId, stripeCustomer)
	if err != nil {
		stripeCustomer, err = customer.Get(user.StripeId, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helpers.ErrorWithCode("server_error", "Failed to get the customer from Stripe"))
			return
		}
		services.GetRedis(c).SetValueForKey(user.StripeId, stripeCustomer)
	}

	stripeCards := []models.Card{}

	for _, paymentSource := range stripeCustomer.Sources.Values {
		if paymentSource.Type == stripe.PaymentSourceCard {
			stripeCard := models.Card{
				Id:       paymentSource.Card.ID,
				Name:     paymentSource.Card.Name,
				Last4:    paymentSource.Card.LastFour,
				Default:  paymentSource.Card.ID == stripeCustomer.DefaultSource.ID,
				ExpMonth: paymentSource.Card.Month,
				ExpYear:  paymentSource.Card.Year,
				Brand:    paymentSource.Card.Brand,
			}

			stripeCards = append(stripeCards, stripeCard)
		}
	}

	c.JSON(http.StatusOK, gin.H{"cards": stripeCards})
}

func (cc CardController) createCustomer(c *gin.Context, user *models.User) (string, error) {
	newCustomer, err := customer.New(&stripe.CustomerParams{
		Email: user.Email,
	})

	if err != nil {
		return "", err
	}

	err = store.UpdateUser(c, params.M{"$set": params.M{"stripeId": newCustomer.ID}})

	if err != nil {
		return "", err
	}

	return newCustomer.ID, nil
}

func (cc CardController) SetDefaultCard(c *gin.Context) {
	user := store.Current(c)

	_, err := customer.Update(
		user.StripeId,
		&stripe.CustomerParams{DefaultSource: c.Param("id")},
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("set_default_card_failed", "Failed to update the customer's default source"))
		return
	}
	services.GetRedis(c).InvalidateObject(user.StripeId)

	c.JSON(http.StatusOK, nil)
}

func (cc CardController) DeleteCard(c *gin.Context) {
	user := store.Current(c)

	_, err := card.Del(
		c.Param("id"),
		&stripe.CardParams{Customer: user.StripeId},
	)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("delete_card_failed", "Failed to delete the customer's card"))
		return
	}

	services.GetRedis(c).InvalidateObject(user.StripeId)

	c.JSON(http.StatusOK, nil)
}
