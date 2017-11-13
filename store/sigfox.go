package store

import (
	"github.com/adrien3d/things-api/models"
	"golang.org/x/net/context"
)

func CreateMessage(c context.Context, message *models.SigfoxMessage) error {
	return FromContext(c).CreateMessage(message)
}

func CreateLocation(c context.Context, location *models.Location) error {
	return FromContext(c).CreateLocation(location)
}
