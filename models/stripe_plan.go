package models

import "github.com/stripe/stripe-go"

type Plan struct {
	Id       string              `json:"id"`
	Amount   uint64              `json:"amount"`
	Interval stripe.PlanInterval `json:"interval"`
	Name     string              `json:"name"`
	Currency stripe.Currency     `json:"currency"`
	MetaData map[string]string   `json:"metadata"` // TODO: MAKE THIS A MODEL TO PREVENT RANDOM DATA
}
