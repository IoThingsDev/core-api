package controllers

import (
	"net/http"

	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
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
}

type Card struct {
	Token string `json:"token" binding:"required"`
}

func NewCardController(mgo *mgo.Database, config *viper.Viper) CardController {
	return CardController{
		mgo,
		config,
	}
}

func (cc CardController) AddCard(c *gin.Context) {
	session := cc.mgo.Session.Copy()
	defer session.Close()
	users := cc.mgo.C(models.UsersCollection)

	userIdInterface, _ := c.Get("userId")
	userId, _ := userIdInterface.(string)
	user := models.User{}

	stripeCard := Card{}
	err := c.Bind(&stripeCard)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	users.FindId(bson.ObjectIdHex(userId)).One(&user)

	if user.StripeId == "" {
		user.StripeId, err = cc.CreateCustomer(&user, users)
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
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("add_card_failed", "Failed to add the card to the customer"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"cards": response})
}

func (cc CardController) GetCards(c *gin.Context) {
	session := cc.mgo.Session.Copy()
	defer session.Close()
	users := cc.mgo.C(models.UsersCollection)

	userIdInterface, _ := c.Get("userId")
	userId, _ := userIdInterface.(string)

	user := models.User{}
	users.FindId(bson.ObjectIdHex(userId)).One(&user)

	if user.StripeId == "" {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("user_not_customer", "The user is not a customer"))
		return
	}

	stripeCustomer, _ := customer.Get(user.StripeId, nil)

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
			}

			stripeCards = append(stripeCards, stripeCard)
		}
	}

	c.JSON(http.StatusOK, gin.H{"cards": stripeCards})
}

func (cc CardController) CreateCustomer(user *models.User, users *mgo.Collection) (string, error) {
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

	return newCustomer.ID, nil
}

func (cc CardController) SetDefaultCard(c *gin.Context) {
	session := cc.mgo.Session.Copy()
	defer session.Close()
	users := cc.mgo.C(models.UsersCollection)

	userIdInterface, _ := c.Get("userId")
	userId, _ := userIdInterface.(string)

	user := models.User{}
	users.FindId(bson.ObjectIdHex(userId)).One(&user)

	stripeCard := Card{}
	err := c.Bind(&stripeCard)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	_, err = customer.Update(
		user.StripeId,
		&stripe.CustomerParams{DefaultSource: stripeCard.Token},
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("set_default_card_failed", "Failed to update the customer's default source"))
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (cc CardController) DeleteCard(c *gin.Context) {
	session := cc.mgo.Session.Copy()
	defer session.Close()
	users := cc.mgo.C(models.UsersCollection)

	userIdInterface, _ := c.Get("userId")
	userId, _ := userIdInterface.(string)

	user := models.User{}
	users.FindId(bson.ObjectIdHex(userId)).One(&user)

	stripeCard := Card{}
	err := c.Bind(&stripeCard)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	_, err = card.Del(
		stripeCard.Token,
		&stripe.CardParams{Customer: user.StripeId},
	)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("delete_card_failed", "Failed to delete the customer's card"))
		return
	}

	c.JSON(http.StatusOK, nil)
}
