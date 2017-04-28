package configuration

import "context"

const (
	storeKey = "config"
)

type Setter interface {
	Set(string, interface{})
}

func FromContext(c context.Context) *config {
	return c.Value(storeKey).(*config)
}

func ToContext(c Setter, config *config) {
	c.Set(storeKey, config)
}

func GetString(c context.Context, key string) string {
	return FromContext(c).GetString(key)
}
