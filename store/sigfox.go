package store

import (
	"github.com/IoThingsDev/api/models"
	"golang.org/x/net/context"
)

// Storing a Sigfox Message
func CreateSigfoxMessage(c context.Context, message *models.SigfoxMessage) error {
	return FromContext(c).CreateMessage(message)
}

// Storing both Sigfox Message and Location
func CreateSigfoxLocationWithMessage(c context.Context, location *models.Location, message *models.SigfoxMessage) error {
	return FromContext(c).CreateLocationWithMessage(location, message)
}

// Get Last Sigfox Messages from Devices
func GetLastDevicesSigfoxMessages(c context.Context) ([]*models.SigfoxMessage, error) {
	return FromContext(c).GetLastDevicesSigfoxMessages(Current(c))
}
