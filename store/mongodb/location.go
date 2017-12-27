package mongodb

import (
	"net/http"

	"github.com/adrien3d/things-api/helpers"
	"github.com/adrien3d/things-api/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *mongo) GetAllDevicesLocations(user *models.User) ([]*models.LastLocation, error) {
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
		{"$sort": bson.M{"location.timestamp": -1}},
		{"$group": bson.M{"_id": "$_id", "name": bson.M{"$first": "$name"}, "location": bson.M{"$push": "$location"}}},
		{"$project": bson.M{"name": "$name", "location": bson.M{"$arrayElemAt": []interface{}{"$location", 0}}}},
	}).All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "get_all_locations_failed", "Failed to get the last locations")
	}

	return list, nil
}
