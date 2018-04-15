package mongodb

import (
	"net/http"

	"github.com/adrien3d/things-api/helpers"
	"github.com/adrien3d/things-api/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *mongo) CreateLocation(location *models.Location) error {
	session := db.Session.Copy()
	defer session.Close()
	locations := db.C(models.LocationsCollection).With(session)

	location.BeforeCreate()
	err := locations.Insert(location)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "location_creation_failed", "Failed to insert the location")
	}

	return nil
}

func (db *mongo) GetLastDevicesLocations(user *models.User) ([]*models.LastLocation, error) {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	list := []*models.LastLocation{}

	err := devices.Pipe([]bson.M{
		{"$match": bson.M{"userId": user.Id}},
		{"$lookup": bson.M{
			"from":         "locations",
			"localField":   "sigfoxId",
			"foreignField": "sigfoxId",
			"as":           "location"}},
		{"$unwind": "$location"},
		{"$sort": bson.M{"location.radius": 1}},
		{"$group": bson.M{"_id": "$_id", "name": bson.M{"$first": "$name"}, "location": bson.M{"$push": "$location"}}},
		{"$project": bson.M{"name": "$name", "location": bson.M{"$arrayElemAt": []interface{}{"$location", 0}}}},
	}).All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "get_all_devices_locations_failed", "Failed to get the last locations")
	}

	return list, nil
}
