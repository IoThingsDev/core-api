package store

import (
	"context"

	"github.com/IoThingsDev/api/helpers/params"
	"github.com/IoThingsDev/api/models"
)

// Create a device in store
func CreateDevice(c context.Context, record *models.Device) error {
	return FromContext(c).CreateDevice(record, Current(c))
}

// Get all devices of a User in store
func GetDevices(c context.Context) ([]*models.Device, error) {
	return FromContext(c).GetDevices(Current(c))
}

// Update a specific device in store
func UpdateDevice(c context.Context, id string, m params.M) error {
	return FromContext(c).UpdateDevice(id, m)
}

// Delete a specific device in store
func DeleteDevice(c context.Context, id string) error {
	return FromContext(c).DeleteDevice(id)
}

// Getting details from a specific device in store
func GetDevice(c context.Context, id string) (*models.Device, error) {
	return FromContext(c).GetDevice(Current(c), id)
}

// Getting last message from a specific device in store
func GetDeviceLastMessage(c context.Context, id string) (*models.SigfoxMessage, error) {
	return FromContext(c).GetDeviceLastMessage(id)
}

// Getting last messages from a specific device in store
func GetLastDeviceMessages(c context.Context, id string) ([]*models.SigfoxMessage, error) {
	return FromContext(c).GetLastDeviceMessages(id)
}

// Getting last locations from a specific device in store
func GetLastDeviceLocations(c context.Context, id string) ([]*models.Location, error) {
	return FromContext(c).GetLastDeviceLocations(id)
}

// Getting all messages from a specific device in store
func GetAllDeviceMessages(c context.Context, id string) ([]*models.SigfoxMessage, error) {
	return FromContext(c).GetAllDeviceMessages(id)
}

// Getting all locations from a specific device in store
func GetAllDeviceLocations(c context.Context, id string) ([]*models.Location, error) {
	return FromContext(c).GetAllDeviceLocations(id)
}
