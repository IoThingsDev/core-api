package store

import (
	"context"

	"github.com/IoThingsDev/api/models"
)

func CreateLocation(c context.Context, location *models.Location) error {
	return FromContext(c).CreateLocation(location)
}
func GetLastDevicesLocations(c context.Context) ([]*models.LastLocation, error) {
	return FromContext(c).GetLastDevicesLocations(Current(c))
}
