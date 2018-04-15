package store

import (
	"context"

	"github.com/IoThingsDev/api/helpers/params"
	"github.com/IoThingsDev/api/models"
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

func GetLastDeviceMessages(c context.Context, id string) ([]*models.SigfoxMessage, error) {
	return FromContext(c).GetLastDeviceMessages(id)
}

func GetLastDeviceLocations(c context.Context, id string) ([]*models.Location, error) {
	return FromContext(c).GetLastDeviceLocations(id)
}

func GetAllDeviceMessages(c context.Context, id string) ([]*models.SigfoxMessage, error) {
	return FromContext(c).GetAllDeviceMessages(id)
}

func GetAllDeviceLocations(c context.Context, id string) ([]*models.Location, error) {
	return FromContext(c).GetAllDeviceLocations(id)
}
