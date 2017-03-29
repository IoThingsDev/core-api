package controllers

import (
	"net/http"

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
}

func NewCardController(mgo *mgo.Database, config *viper.Viper) *CardController {
	return &CardController{
		mgo,
		config,
	}
}

func (cc CardController) AddCard(c *gin.Context) {
	session := cc.mgo.Session.Copy()
	defer session.Close()
	users := cc.mgo.C(models.UsersCollection)

	services.SetStripeKeyAndBackend(cc.config)

	userIdInterface, _ := c.Get("userId")
	userId, _ := userIdInterface.(string)
	user := models.User{}

	stripeCard := stripe.Card{}
	err := c.Bind(&stripeCard)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	users.FindId(bson.ObjectIdHex(userId)).One(&user)

	if user.StripeId == "" {
		user.StripeId, err = cc.CreateCustomer(&user, users)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err) //helpers.ErrorWithCode("server_error", "Failed to create the customer in our billing platform"))
			return
		}
	}

	_, err = card.New(&stripe.CardParams{
		Customer: user.StripeId,
		Token:    stripeCard.ID,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err) //helpers.ErrorWithCode("add_card_failed", "Failed to add the card to the customer"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Card successfully added to the user"})
}

func (cc CardController) GetCards(c *gin.Context) {
	session := cc.mgo.Session.Copy()
	defer session.Close()
	users := cc.mgo.C(models.UsersCollection)

	services.SetStripeKeyAndBackend(cc.config)

	userIdInterface, _ := c.Get("userId")
	userId, _ := userIdInterface.(string)

	user := models.User{}
	users.FindId(bson.ObjectIdHex(userId)).One(&user)

	if user.StripeId == "" {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("user_not_customer", "The user is not a customer"))
		return
	}

	stripeCards := []*stripe.Card{}
	params := &stripe.CardListParams{Customer: user.StripeId}
	i := card.List(params)
	for i.Next() {
		stripeCards = append(stripeCards, i.Card())
	}

	c.JSON(http.StatusCreated, gin.H{"cards": stripeCards})
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
