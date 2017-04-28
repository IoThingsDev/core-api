package repositories

import "github.com/dernise/base-api/models"

type Store interface {
	CreateUser(*models.User) error
	GetUser(string) (*models.User, error)
	ActivateUser(string, string) error
}
