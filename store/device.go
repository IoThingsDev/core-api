package store

import (
	"context"

	"github.com/adrien3d/things-api/helpers/params"
	"github.com/adrien3d/things-api/models"
)

func CreateDevice(c context.Context, record *models.Device) error {
	return FromContext(c).CreateDevice(record, Current(c))
}

func GetDevices(c context.Context) ([]*models.Device, error) {
	return FromContext(c).GetDevices(Current(c))
}

func UpdateDevice(c context.Context, id string, m params.M) error {
	return FromContext(c).UpdateDevice(id, m)
}

func DeleteDevice(c context.Context, id string) error {
	return FromContext(c).DeleteDevice(id)
}

func GetDevice(c context.Context, id string) (*models.Device, error) {
	return FromContext(c).GetDevice(Current(c), id)
}

func GetLastMessages(c context.Context, id string) ([]*models.SigfoxMessage, error) {
	return FromContext(c).GetLastMessages(id)
}

func GetLastLocations(c context.Context, id string) ([]*models.Location, error) {
	return FromContext(c).GetLastLocations(id)
}
