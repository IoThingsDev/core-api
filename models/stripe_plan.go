package models

import "github.com/stripe/stripe-go"

type Plan struct {
	Id       string              `json:"id" binding:"required"`
	Amount   uint64              `json:"amount"`
	Interval stripe.PlanInterval `json:"interval" binding:"required"`
	Name     string              `json:"name" binding:"required"`
	Currency stripe.Currency     `json:"currency" binding:"required"`
	MetaData map[string]string   `json:"metadata" binding:"required"` // TODO: MAKE THIS A MODEL TO PREVENT RANDOM DATA
}
