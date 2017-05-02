package store

import "github.com/dernise/base-api/models"

type Store interface {
	CreateUser(*models.User) error
	FindUserById(string) (*models.User, error)
	ActivateUser(string, string) error
	FindUser(map[string]interface{}) (*models.User, error)
}
