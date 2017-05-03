package store

import (
	"github.com/dernise/base-api/models"
	"golang.org/x/net/context"
)

const (
	currentKey = "currentUser"
	storeKey   = "store"
)

type Setter interface {
	Set(string, interface{})
}

func Current(c context.Context) *models.User {
	return c.Value(currentKey).(*models.User)
}

func FromContext(c context.Context) Store {
	return c.Value(storeKey).(Store)
}

func ToContext(c Setter, store Store) {
	c.Set(storeKey, store)
}
