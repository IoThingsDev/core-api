package store

import (
	"github.com/adrien3d/things-api/models"
	"golang.org/x/net/context"
)

const (
	CurrentKey    = "currentUser"
	LoginTokenKey = "currentLoginToken"
	StoreKey      = "store"
)

type Setter interface {
	Set(string, interface{})
}

func Current(c context.Context) *models.User {
	return c.Value(CurrentKey).(*models.User)
}

func ToContext(c Setter, store Store) {
	c.Set(StoreKey, store)
}

func FromContext(c context.Context) Store {
	return c.Value(StoreKey).(Store)
}
