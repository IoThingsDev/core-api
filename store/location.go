package store

import (
	"context"

	"github.com/adrien3d/things-api/models"
)

func GetAllDevicesLocations(c context.Context) ([]*models.LastLocation, error) {
	return FromContext(c).GetAllDevicesLocations(Current(c))
}
