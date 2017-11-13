package store

import (
	"context"

	"github.com/adrien3d/things-api/helpers/params"
	"github.com/adrien3d/things-api/models"
)

func CreateUser(c context.Context, record *models.User) error {
	return FromContext(c).CreateUser(record)
}

func FindUserById(c context.Context, id string) (*models.User, error) {
	return FromContext(c).FindUserById(id)
}

func FindUser(c context.Context, params params.M) (*models.User, error) {
	return FromContext(c).FindUser(params)
}

func ActivateUser(c context.Context, activationKey string, id string) error {
	return FromContext(c).ActivateUser(activationKey, id)
}

func UpdateUser(c context.Context, params params.M) error {
	return FromContext(c).UpdateUser(Current(c), params)
}

func AddLoginToken(c context.Context, user *models.User, ip string) (*models.LoginToken, error) {
	return FromContext(c).AddLoginToken(user, ip)
}

func RemoveLoginToken(c context.Context) error {
	return FromContext(c).RemoveLoginToken(Current(c), c.Value(LoginTokenKey).(string))
}
