package store

import (
	"context"

	"github.com/IoThingsDev/api/models"
)

// Create a single Location in store
func CreateLocation(c context.Context, location *models.Location) error {
	return FromContext(c).CreateLocation(location)
}

// Getting Last Locations from all devices of a user in store
func GetLastDevicesLocations(c context.Context) ([]*models.LastLocation, error) {
	return FromContext(c).GetLastDevicesLocations(Current(c))
}
