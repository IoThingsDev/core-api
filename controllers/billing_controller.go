package controllers

import (
	"net/http"
	"time"

	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
	"github.com/dernise/base-api/services"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/currency"
	"github.com/stripe/stripe-go/plan"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
)

type BillingController struct {
	mgo         *mgo.Database
	emailSender services.EmailSender
	config      *viper.Viper
	redis       *services.Redis
}

func NewBillingController(mgo *mgo.Database, emailSender services.EmailSender, config *viper.Viper, redis *services.Redis) BillingController {
	return BillingController{
		mgo,
		emailSender,
		config,
		redis,
	}
}

func (bc BillingController) CreateTransaction(c *gin.Context) {
	session := bc.mgo.Session.Copy()
	defer session.Close()
	transactions := bc.mgo.C(models.TransactionsCollection).With(session)
	users := bc.mgo.C(models.UsersCollection).With(session)

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
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("no_card_found", "The customer doesn't have any card to pay"))
		return
	}

	chargeParams := &stripe.ChargeParams{
		Amount:   transaction.Amount,
		Currency: currency.EUR,
		Customer: user.StripeId,
	}
	response, err := charge.New(chargeParams)
	if err != nil {
		transaction.Error = err.Error()
		transaction.Failed = false
	} else if response.Status != "succeeded" {
		transaction.Error = response.FailCode
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

func (bc BillingController) GetPlans(c *gin.Context) {
	stripePlans := []models.Plan{}
	err := bc.redis.GetValueForKey("billing-plans", &stripePlans)
	if err != nil {
		i := plan.List(nil)

		for i.Next() {
			p := i.Plan()
			stripePlan := models.Plan{
				Id:       p.ID,
				Amount:   p.Amount,
				Currency: p.Currency,
				Interval: p.Interval,
				Name:     p.Name,
				MetaData: p.Meta,
			}

			stripePlans = append(stripePlans, stripePlan)
		}

		bc.redis.SetValueForKey("billing-plans", stripePlans)
	}

	c.JSON(http.StatusCreated, gin.H{"plans": stripePlans})
}

func (bc BillingController) CreatePlan(c *gin.Context) {
	stripePlan := models.Plan{}
	err := c.Bind(&stripePlan)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	params := stripe.Params{
		Meta: stripePlan.MetaData,
	}

	_, err = plan.New(&stripe.PlanParams{
		Amount:   stripePlan.Amount,
		Interval: stripePlan.Interval,
		Name:     stripePlan.Name,
		Currency: stripePlan.Currency,
		ID:       stripePlan.Id,
		Params:   params,
	})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("plan_creation_failed", "The plan has not been created"))
		return
	}

	bc.redis.InvalidateObject("billing-plans")

	c.JSON(http.StatusCreated, gin.H{"plans": stripePlan})
}
