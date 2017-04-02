package controllers

import (
	"net/http"

	"fmt"
	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
	"github.com/dernise/base-api/services"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/customer"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CardController struct {
	mgo    *mgo.Database
	config *viper.Viper
	redis  *services.Redis
}

type Card struct {
	Token string `json:"token" binding:"required"`
}

func NewCardController(mgo *mgo.Database, config *viper.Viper, redis *services.Redis) CardController {
	return CardController{
		mgo,
		config,
		redis,
	}
}

func (cc CardController) AddCard(c *gin.Context) {
	userInterface, _ := c.Get("currentUser")
	user := userInterface.(models.User)

	stripeCard := Card{}
	err := c.Bind(&stripeCard)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	if user.StripeId == "" {
		user.StripeId, err = cc.CreateCustomer(&user)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("server_error", "Failed to create the customer in our billing platform"))
			return
		}
	}

	response, err := card.New(&stripe.CardParams{
		Customer: user.StripeId,
		Token:    stripeCard.Token,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("add_card_failed", "Failed to add the card to the customer"))
		return
	}
	cc.redis.InvalidateObject(user.StripeId)

	c.JSON(http.StatusCreated, gin.H{"cards": response})
}

func (cc CardController) GetCards(c *gin.Context) {
	userInterface, _ := c.Get("currentUser")
	user := userInterface.(models.User)

	if user.StripeId == "" {
		c.JSON(http.StatusOK, gin.H{"cards": []models.Card{}})
		return
	}

	stripeCustomer := &stripe.Customer{}

	err := cc.redis.GetValueForKey(user.StripeId, stripeCustomer)
	if err != nil {
		stripeCustomer, err = customer.Get(user.StripeId, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helpers.ErrorWithCode("server_error", "Failed to get the customer from Stripe"))
			return
		}
		cc.redis.SetValueForKey(user.StripeId, stripeCustomer)
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

func (cc CardController) CreateCustomer(user *models.User) (string, error) {
	session := cc.mgo.Session.Copy()
	defer session.Close()
	users := cc.mgo.C(models.UsersCollection)

	customerParams := &stripe.CustomerParams{
		Email: user.Email,
	}

	newCustomer, err := customer.New(customerParams)
	if err != nil {
		return "", err
	}

	err = users.UpdateId(user.Id, bson.M{"$set": bson.M{"stripeId": newCustomer.ID}})
	if err != nil {
		return "", err
	}
	cc.redis.InvalidateObject(user.Id.Hex())

	return newCustomer.ID, nil
}

func (cc CardController) SetDefaultCard(c *gin.Context) {
	userInterface, _ := c.Get("currentUser")
	user := userInterface.(models.User)

	_, err := customer.Update(
		user.StripeId,
		&stripe.CustomerParams{DefaultSource: c.Param("id")},
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("set_default_card_failed", "Failed to update the customer's default source"))
		return
	}
	cc.redis.InvalidateObject(user.StripeId)

	c.JSON(http.StatusOK, nil)
}

func (cc CardController) DeleteCard(c *gin.Context) {
	userInterface, _ := c.Get("currentUser")
	user := userInterface.(models.User)

	_, err := card.Del(
		c.Param("id"),
		&stripe.CardParams{Customer: user.StripeId},
	)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("delete_card_failed", "Failed to delete the customer's card"))
		return
	}
	cc.redis.InvalidateObject(user.StripeId)

	c.JSON(http.StatusOK, nil)
}
