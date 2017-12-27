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

	CreateDevice(*models.Device, *models.User) error
	GetDevices(*models.User) ([]*models.Device, error)
	UpdateDevice(string, params.M) error
	DeleteDevice(string) error
	GetDevice(*models.User, string) (*models.Device, error)
	GetLastMessages(string) ([]*models.SigfoxMessage, error)
	GetLastLocations(string) ([]*models.Location, error)
	GetAllMessages(string) ([]*models.SigfoxMessage, error)
	GetAllLocations(string) ([]*models.Location, error)
	GetAllDevicesLocations(*models.User) ([]*models.LastLocation, error)

	AddLoginToken(*models.User, string) (*models.LoginToken, error)
	RemoveLoginToken(*models.User, string) error
}
