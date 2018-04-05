package store

import (
	"github.com/adrien3d/things-api/models"
	"golang.org/x/net/context"
)

func CreateSigfoxMessage(c context.Context, message *models.SigfoxMessage) error {
	return FromContext(c).CreateMessage(message)
}

func CreateSigfoxLocationWithMessage(c context.Context, location *models.Location, message *models.SigfoxMessage) error {
	return FromContext(c).CreateLocationWithMessage(location, message)
}

func GetLastDevicesSigfoxMessages(c context.Context) ([]*models.LastLocation, error) {
	return FromContext(c).GetLastDevicesSigfoxMessages(Current(c))
}
