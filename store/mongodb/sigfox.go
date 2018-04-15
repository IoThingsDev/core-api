package mongodb

import (
	"net/http"

	"github.com/IoThingsDev/api/helpers"
	"github.com/IoThingsDev/api/helpers/params"
	"github.com/IoThingsDev/api/models"
	"gopkg.in/mgo.v2/bson"
)

func (db *mongo) CreateMessage(message *models.SigfoxMessage) error {
	session := db.Session.Copy()
	defer session.Close()
	sigfoxMessages := db.C(models.SigfoxMessagesCollection).With(session)

	message.BeforeCreate()
	err := sigfoxMessages.Insert(message)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "message_creation_failed", "Failed to insert the sigfox message")
	}

	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	err = devices.Find(params.M{"sigfoxId": message.SigfoxId}).One(device)
	if err != nil {
		return helpers.NewError(http.StatusPartialContent, "sigfox_device_id_not_found", "Device Sigfox ID not found")
	} else {
		err = devices.Update(bson.M{"sigfoxId": message.SigfoxId}, bson.M{"$set": bson.M{"lastAcc": message.Timestamp}})
		if err != nil {
			return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update device last activity")
		}

		err = devices.Update(bson.M{"sigfoxId": message.SigfoxId}, bson.M{"$set": bson.M{"active": true}})
		if err != nil {
			return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update device status")
		}
	}

	return nil
}

func (db *mongo) CreateLocationWithMessage(location *models.Location, message *models.SigfoxMessage) error {
	session := db.Session.Copy()
	defer session.Close()
	locations := db.C(models.LocationsCollection).With(session)

	location.BeforeCreate()
	err := locations.Insert(location)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "location_creation_failed", "Failed to insert the location")
	}

	sigfoxMessages := db.C(models.SigfoxMessagesCollection).With(session)

	message.BeforeCreate()
	err = sigfoxMessages.Insert(message)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "message_creation_failed", "Failed to insert the sigfox message")
	}

	devices := db.C(models.DevicesCollection).With(session)

	devices.Update(bson.M{"sigfoxId": message.SigfoxId}, bson.M{"$set": bson.M{"lastAcc": message.Timestamp}})

	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update the device")
	}

	return nil
}

func (db *mongo) GetLastDevicesSigfoxMessages(user *models.User) ([]*models.SigfoxMessage, error) {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	list := []*models.SigfoxMessage{}

	err := devices.Pipe([]bson.M{
		{"$match": bson.M{"userId": user.Id}},
		{"$lookup": bson.M{
			"from":         "sigfox_messages",
			"localField":   "sigfoxId",
			"foreignField": "sigfoxId",
			"as":           "message"}},
		{"$unwind": "message"},
		//{"$sort": bson.M{"location.radius": 1}},
		{"$group": bson.M{"_id": "$_id", "name": bson.M{"$first": "$name"}, "message": bson.M{"$push": "message"}}}, //rendered
		{"$project": bson.M{"name": "$name", "message": bson.M{"$arrayElemAt": []interface{}{"message", 0}}}},
	}).All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "get_last_messages_failed", "Failed to get the last messages")
	}

	return list, nil
}
