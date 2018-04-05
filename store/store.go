package store

import (
	"github.com/adrien3d/things-api/helpers/params"
	"github.com/adrien3d/things-api/models"
)

type Store interface {
	CreateUser(*models.User) error
	FindUserById(string) (*models.User, error)
	ActivateUser(string, string) error
	FindUser(params.M) (*models.User, error)
	UpdateUser(*models.User, params.M) error

	CreateMessage(*models.SigfoxMessage) error
	CreateLocation(*models.Location) error
	CreateLocationWithMessage(*models.Location, *models.SigfoxMessage) error

	GetLastDevicesSigfoxMessages(*models.User) ([]*models.LastLocation, error)

	CreateDevice(*models.Device, *models.User) error
	GetDevices(*models.User) ([]*models.Device, error)
	UpdateDevice(string, params.M) error
	DeleteDevice(string) error
	GetDevice(*models.User, string) (*models.Device, error)
	GetLastDeviceMessages(string) ([]*models.SigfoxMessage, error)
	GetLastDeviceLocations(string) ([]*models.Location, error)
	GetAllDeviceMessages(string) ([]*models.SigfoxMessage, error)
	GetAllDeviceLocations(string) ([]*models.Location, error)

	GetLastDevicesLocations(*models.User) ([]*models.LastLocation, error)

	AddLoginToken(*models.User, string) (*models.LoginToken, error)
	RemoveLoginToken(*models.User, string) error
}
