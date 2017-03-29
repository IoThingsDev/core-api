package services

import (
	"io"

	"github.com/spf13/viper"
	"github.com/stripe/stripe-go"
)

type FakeStripeBackend struct{}

func (fsb FakeStripeBackend) Call(method, path, key string, body *stripe.RequestValues, params *stripe.Params, v interface{}) error {
	charge := v.(*stripe.Charge)
	charge.Status = "succeeded"
	return nil
}

func (fsb FakeStripeBackend) CallMultipart(method, path, key, boundary string, body io.Reader, params *stripe.Params, v interface{}) error {
	return nil
}

func SetStripeKeyAndBackend(config *viper.Viper) {
	stripe.Key = config.GetString("stripe_api_key")

	if config.GetString("env") == "testing" {
		backend := FakeStripeBackend{}
		stripe.SetBackend("api", backend)
	}
}
