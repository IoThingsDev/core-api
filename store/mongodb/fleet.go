package mongodb

import (
	"github.com/IoThingsDev/api/helpers"
	"github.com/IoThingsDev/api/helpers/params"
	"github.com/IoThingsDev/api/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func (db *mongo) GetAllFleets() ([]models.Fleet, error) {
	session := db.Session.Copy()
	defer session.Close()
	fleetCollection := db.C(models.FleetsCollection).With(session)

	fleets := []models.Fleet{}
	err := fleetCollection.Find(bson.M{}).All(&fleets)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_fleets_failed", "Failed to get the fleets: "+err.Error())
	}

	return fleets, nil
}

func (db *mongo) CreateFleet(fleet *models.Fleet) error {
	session := db.Session.Copy()
	defer session.Close()
	fleetCollection := db.C(models.FleetsCollection).With(session)

	fleet.Id = bson.NewObjectId().Hex()

	err := fleetCollection.Insert(fleet)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "fleet_creation_failed", "Failed to insert the fleet in the database: "+err.Error())
	}

	return nil
}

func (db *mongo) GetFleetById(id string) (*models.Fleet, error) {
	session := db.Session.Copy()
	defer session.Close()
	fleetCollection := db.C(models.FleetsCollection).With(session)

	fleet := &models.Fleet{}
	err := fleetCollection.FindId(id).One(fleet)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "fleet_not_found", "Could not find the fleet")
	}

	return fleet, nil
}

func (db *mongo) UpdateFleet(id string, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	fleets := db.C(models.FleetsCollection).With(session)

	err := fleets.UpdateId(id, params)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "fleet_update_failed", "Failed to update the fleets: "+err.Error())
	}

	return nil
}

func (db *mongo) GetLastFleetMessages(id string) ([]*models.SigfoxMessage, error) {
	list := []*models.SigfoxMessage{}

	fleet, err := db.GetFleetById(id)
	if err != nil {
		return nil, err
	}

	for _, id := range fleet.DeviceIds {
		lastMessage, err := db.GetDeviceLastMessage(id)
		if err != nil {
			continue
		}

		list = append(list, lastMessage)
	}

	return list, nil
}

func (db *mongo) DeleteFleet(id string) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.FleetsCollection).With(session)

	err := devices.RemoveId(id)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_delete_failed", "Failed to delete the device")
	}

	return nil
}
