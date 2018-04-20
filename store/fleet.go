package store

import (
	"context"
	"github.com/IoThingsDev/api/models"
)

func GetAllFleets(c context.Context) ([]models.Fleet, error) {
	return FromContext(c).GetAllFleets()
}

func CreateFleet(c context.Context, fleet *models.Fleet) error {
	return FromContext(c).CreateFleet(fleet)
}

func GetFleetById(c context.Context, id string) (*models.Fleet, error) {
	return FromContext(c).GetFleetById(id)
}

func GetLastFleetMessages(c context.Context, id string) ([]*models.SigfoxMessage, error) {
	return FromContext(c).GetLastFleetMessages(id)
}

func DeleteFleet(c context.Context, id string) error {
	return FromContext(c).DeleteFleet(id)
}
