package controllers

import (
	"net/http"

	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
	"github.com/dernise/base-api/services"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/plan"
	"gopkg.in/gin-gonic/gin.v1"
)

type BillingController struct{}

func NewBillingController() BillingController {
	return BillingController{}
}

func (bc BillingController) GetPlans(c *gin.Context) {
	stripePlans := []models.Plan{}
	err := services.GetRedis(c).GetValueForKey("billing-plans", &stripePlans)
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

		services.GetRedis(c).SetValueForKey("billing-plans", stripePlans)
	}

	c.JSON(http.StatusOK, gin.H{"plans": stripePlans})
}

func (bc BillingController) CreatePlan(c *gin.Context) {
	stripePlan := models.Plan{}
	err := c.BindJSON(&stripePlan)
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

	services.GetRedis(c).InvalidateObject("billing-plans")

	c.JSON(http.StatusCreated, gin.H{"plans": stripePlan})
}
