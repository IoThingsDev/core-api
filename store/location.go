package store

import (
	"context"

	"github.com/adrien3d/things-api/models"
)

func GetAllLocations(c context.Context) ([]*models.LastLocation, error) {
	return FromContext(c).GetAllLocations(Current(c))
}
