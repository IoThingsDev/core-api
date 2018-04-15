package store

import (
	"github.com/IoThingsDev/api/models"
	"golang.org/x/net/context"
)

func CreateSigfoxMessage(c context.Context, message *models.SigfoxMessage) error {
	return FromContext(c).CreateMessage(message)
}

func CreateSigfoxLocationWithMessage(c context.Context, location *models.Location, message *models.SigfoxMessage) error {
	return FromContext(c).CreateLocationWithMessage(location, message)
}

func GetLastDevicesSigfoxMessages(c context.Context) ([]*models.SigfoxMessage, error) {
	return FromContext(c).GetLastDevicesSigfoxMessages(Current(c))
}
