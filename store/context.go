package store

import (
	"github.com/IoThingsDev/api/models"
	"golang.org/x/net/context"
)

//Keys for various values
const (
	CurrentKey    = "currentUser"
	LoginTokenKey = "currentLoginToken"
	StoreKey      = "store"
)

//Sets a value
type Setter interface {
	Set(string, interface{})
}

//Get Current context
func Current(c context.Context) *models.User {
	return c.Value(CurrentKey).(*models.User)
}

//Set a value in store
func ToContext(c Setter, store Store) {
	c.Set(StoreKey, store)
}

//Get value from Store
func FromContext(c context.Context) Store {
	return c.Value(StoreKey).(Store)
}
