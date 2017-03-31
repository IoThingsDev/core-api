package services

import (
	"io"

	"github.com/spf13/viper"
	"github.com/stripe/stripe-go"
)

type FakeStripeBackend struct{}

func (fsb FakeStripeBackend) Call(method, path, key string, body *stripe.RequestValues, params *stripe.Params, v interface{}) error {
	if charge, ok := v.(*stripe.Charge); ok {
		charge.Status = "succeeded"
	} else if customer, ok := v.(*stripe.Customer); ok {
		customer.Sources = &stripe.SourceList{}

		customer.Sources.Values = append(customer.Sources.Values, &stripe.PaymentSource{
			Type: stripe.PaymentSourceCard,
			ID:   "testId",
			Card: &stripe.Card{},
		})

		customer.DefaultSource = &stripe.PaymentSource{ID: "testId"}
	}

	return nil
}

func (fsb FakeStripeBackend) CallMultipart(method, path, key, boundary string, body io.Reader, params *stripe.Params, v interface{}) error {
	return nil
}

func SetStripeKeyAndBackend(config *viper.Viper) {
	stripe.Key = config.GetString("stripe_api_key")

	if config.GetString("env") == "testing" {
		stripe.SetBackend("api", FakeStripeBackend{})
	}
}
