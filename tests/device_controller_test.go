package tests

import (
	"net/http"
	"testing"

	"github.com/IoThingsDev/api/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestDeviceController(t *testing.T) {
	deviceReq := []byte(`
	{
		"name":"Device 1",
		"sigfoxId":"AABBCC",
		"bleMac": "AA:BB:CC:DD:EE",
		"lastAcc": 1491209490,
		"active": true
	}`)

	resp := SendRequestWithToken(deviceReq, "POST", "/v1/devices/", authToken)
	assert.Equal(t, http.StatusCreated, resp.Code)

	deviceReq = []byte(`
	{
		"name":"Device 1",
		"sigfoxId":"AABBCC",
		"bleMac": "AA:BB:CC:DD:EO",
		"lastAcc": 1491209490,
		"active": true
	}`)

	resp = SendRequestWithToken(deviceReq, "POST", "/v1/devices/", authToken)
	assert.Equal(t, http.StatusConflict, resp.Code)

	// Test activation
	device := models.Device{}
	err := api.Database.C(models.DevicesCollection).Find(bson.M{"sigfoxId": "17CA54"}).One(&device)
	if err != nil {
		t.Fail()
		return
	}

	resp = SendRequestWithToken(nil, "GET", "/v1/devices/"+device.Id, authToken)
	assert.Equal(t, http.StatusOK, resp.Code)

	deviceReq = []byte(`
	{
		"name":"Device 1",
		"sigfoxId":"AABBCC",
		"bleMac": "AA:3B:CC:DD:EE",
		"lastAcc": 1491809490,
		"active": false
	}`)

	resp = SendRequestWithToken(deviceReq, "PUT", "/v1/devices/"+device.Id, authToken)
	assert.Equal(t, http.StatusOK, resp.Code)

	assert.Equal(t, device.Active, true)

	resp = SendRequestWithToken(nil, "GET", "/v1/devices/", authToken)
	assert.Equal(t, http.StatusOK, resp.Code)

	resp = SendRequestWithToken(nil, "DELETE", "/v1/devices/"+device.Id, authToken)
	assert.Equal(t, http.StatusOK, resp.Code)
}
