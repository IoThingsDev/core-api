package controllers

import (
	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
	"github.com/dernise/base-api/services"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/currency"
	"github.com/stripe/stripe-go/customer"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type BillingController struct {
	mgo         *mgo.Database
	emailSender services.EmailSender
	config      *viper.Viper
}

func NewBillingController(mgo *mgo.Database, emailSender services.EmailSender, config *viper.Viper) *BillingController {
	return &BillingController{
		mgo,
		emailSender,
		config,
	}
}

func (bc BillingController) CreateTransaction(c *gin.Context) {
	session := bc.mgo.Session.Copy()
	defer session.Close()
	transactions := bc.mgo.C(models.TransactionsCollection).With(session)
	users := bc.mgo.C(models.UsersCollection).With(session)

	stripe.Key = bc.config.GetString("stripe_api_key")

	transaction := models.Transaction{}
	err := c.Bind(&transaction)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	user := models.User{}
	if err = users.FindId(transaction.UserId).One(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to find the user for the transaction"))
		return
	}

	if user.StripeId == "" {
		user.StripeId, err = bc.CreateCustomer(&user, users, transaction.CardToken)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err) //helpers.ErrorWithCode("server_error", "Failed to create the customer in our billing platform"))
			return
		}
	}

	chargeParams := &stripe.ChargeParams{
		Amount:   transaction.Amount,
		Currency: currency.EUR,
		Customer: user.StripeId,
	}
	charge, err := charge.New(chargeParams)
	if err != nil {
		transaction.Error = err.Error()
		transaction.Failed = false
	} else if charge.Status != "succeeded" {
		transaction.Error = charge.FailCode
		transaction.Failed = false
	} else {
		transaction.Failed = true
	}

	transaction.Date = time.Now()
	transactions.Insert(transaction)

	if transaction.Failed == false {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("payment_failed", "The payment failed"))
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Payment successed"})
	}
}

func (bc BillingController) CreateCustomer(user *models.User, users *mgo.Collection, token string) (string, error) {
	customerParams := &stripe.CustomerParams{
		Email: user.Email,
	}

	customerParams.SetSource(token)

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
