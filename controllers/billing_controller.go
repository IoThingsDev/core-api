package controllers

import (
	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
	"github.com/dernise/base-api/services"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"os"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/currency"
	"time"
)

type BillingController struct {
	mgo         *mgo.Database
	emailSender services.EmailSender
}

func NewBillingController(mgo *mgo.Database, emailSender services.EmailSender) *BillingController {
	return &BillingController{
		mgo,
		emailSender,
	}
}

func (bc BillingController) CreateTransaction(c *gin.Context) {
	session := bc.mgo.Session.Copy()
	defer session.Close()
	transactions := bc.mgo.C(models.TransactionsCollection).With(session)
	users := bc.mgo.C(models.UsersCollection).With(session)

	stripe.Key = os.Getenv("STRIPE_API_KEY")

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
			c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("server_error", "Failed to create the customer in our billing platform"))
		}
	}

	chargeParams := &stripe.ChargeParams{
		Amount: transaction.Amount,
		Currency: currency.EUR,
		Customer: user.StripeId,
	}
	charge, err := charge.New(chargeParams)
	if err != nil {
		transaction.Error = err.Error()
		transaction.Status = false
	} else if charge.Status != "succeeded" {
		transaction.Error = charge.FailCode
		transaction.Status = false
	} else {
		transaction.Status = true
	}

	transaction.Date = time.Now()
	transactions.Insert(transaction)

	if transaction.Status == false {
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
