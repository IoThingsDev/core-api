package mongodb

import (
	"net/http"

	"github.com/IoThingsDev/api/helpers"
	"github.com/IoThingsDev/api/helpers/params"
	"github.com/IoThingsDev/api/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB: Create a device
func (db *mongo) CreateDevice(device *models.Device, user *models.User) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	device.UserId = user.Id
	device.BeforeCreate()

	if device.SigfoxId != "" {
		count, _ := devices.Find(params.M{"sigfoxId": device.SigfoxId}).Count()
		if count > 0 {
			return helpers.NewError(http.StatusConflict, "device_already_exists", "Device already exists")
		}
	}

	if device.BLEMac != "" {
		count, _ := devices.Find(params.M{"bleMac": device.BLEMac}).Count()
		if count > 0 {
			return helpers.NewError(http.StatusConflict, "device_already_exists", "Device already exists")
		}
	}

	err := devices.Insert(device)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_creation_failed", "Failed to create the device")
	}

	return nil
}

// MongoDB: Get all devices of a User
func (db *mongo) GetDevices(user *models.User) ([]*models.Device, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)

	list := []*models.Device{}
	err := devices.Find(params.M{"userId": user.Id}).Sort("-lastAcc").All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "devices_not_found", "Devices not found")
	}

	return list, nil
}

// MongoDB: Update a specific device
func (db *mongo) UpdateDevice(id string, m params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	change := mgo.Change{
		Update:    m,
		Upsert:    false,
		Remove:    false,
		ReturnNew: false,
	}
	_, err := devices.Find(bson.M{"_id": id}).Apply(change, nil)

	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_update_failed", "Failed to update the device")
	}

	return nil
}

// MongoDB: Delete a specific device
func (db *mongo) DeleteDevice(id string) error {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)

	err := devices.RemoveId(id)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "device_delete_failed", "Failed to delete the device")
	}

	return nil
}

// MongoDB: Getting details from a specific device
func (db *mongo) GetDevice(user *models.User, id string) (*models.Device, error) {
	session := db.Session.Copy()
	defer session.Close()

	devices := db.C(models.DevicesCollection).With(session)
	device := &models.Device{}

	err := devices.FindId(id).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "device_not_found", "Device not found")
	}

	return device, nil
}

// MongoDB: Getting last messages from a specific device
func (db *mongo) GetLastDeviceMessages(id string) ([]*models.SigfoxMessage, error) {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)
	sigfoxMessages := db.C(models.SigfoxMessagesCollection).With(session)

	device := &models.Device{}

	err := devices.FindId(id).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_failed", "Failed to find the device")
	}

	list := []*models.SigfoxMessage{}
	err = sigfoxMessages.Find(bson.M{"sigfoxId": device.SigfoxId}).Limit(10).Sort("-$natural").All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_failed", "Failed to query the Database")
	}

	return list, nil
}

// MongoDB: Getting last message from a specific device
func (db *mongo) GetDeviceLastMessage(id string) (*models.SigfoxMessage, error) {
	session := db.Session.Copy()
	defer session.Close()
	sigfoxMessages := db.C(models.SigfoxMessagesCollection).With(session)

	var message models.SigfoxMessage
	err := sigfoxMessages.Find(bson.M{"sigfoxId": id}).Sort("-$natural").One(&message)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_failed", "Failed to query the Database: "+err.Error())
	}

	return &message, nil
}

// MongoDB: Getting last locations from a specific device
func (db *mongo) GetLastDeviceLocations(id string) ([]*models.Location, error) {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)
	locations := db.C(models.LocationsCollection).With(session)

	device := &models.Device{}

	err := devices.FindId(id).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_failed", "Failed to find the device")
	}

	list := []*models.Location{}
	err = locations.Find(bson.M{"sigfoxId": device.SigfoxId}).Limit(10).Sort("-$natural").All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_failed", "Failed to query the Database")
	}

	return list, nil
}

// MongoDB: Getting all messages from a specific device
func (db *mongo) GetAllDeviceMessages(id string) ([]*models.SigfoxMessage, error) {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)
	sigfoxMessages := db.C(models.SigfoxMessagesCollection).With(session)

	device := &models.Device{}

	err := devices.FindId(id).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_failed", "Failed to find the device")
	}

	list := []*models.SigfoxMessage{}
	err = sigfoxMessages.Find(bson.M{"sigfoxId": device.SigfoxId}).Limit(100).Sort("-$natural").All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_failed", "Failed to query the Database")
	}

	return list, nil
}

// MongoDB: Getting all locations from a specific device
func (db *mongo) GetAllDeviceLocations(id string) ([]*models.Location, error) {
	session := db.Session.Copy()
	defer session.Close()
	devices := db.C(models.DevicesCollection).With(session)
	locations := db.C(models.LocationsCollection).With(session)

	device := &models.Device{}

	err := devices.FindId(id).One(device)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_failed", "Failed to find the device")
	}

	list := []*models.Location{}
	err = locations.Find(bson.M{"sigfoxId": device.SigfoxId}).Limit(100).Sort("-$natural").All(&list)
	if err != nil {
		return nil, helpers.NewError(http.StatusInternalServerError, "query_failed", "Failed to query the Database")
	}

	return list, nil
}
