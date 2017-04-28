package repositories

import (
	"context"

	"github.com/dernise/base-api/models"
)

func CreateUser(c context.Context, record *models.User) error {
	return FromContext(c).CreateUser(record)
}

func GetUser(c context.Context, id string) (*models.User, error) {
	return FromContext(c).GetUser(id)
}

func ActivateUser(c context.Context, activationKey string, id string) error {
	return FromContext(c).ActivateUser(activationKey, id)
}
