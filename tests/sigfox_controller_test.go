package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMessage(t *testing.T) {
	parameters := []byte(`
	{
		"sigfoxId":"17CA54",
		"frameNumber":2,
		"timestamp": 1491209486,
		"station": "12C2",
		"snr": 15.78,
		"avgSnr": 12.93,
		"rssi": 27.92,
		"latGps": 50.4350385,
		"lngGps": 2.82355960,
		"radiusGps": 3.98,
		"latSf": 50.4308385,
		"lngSf": 2.82385960,
		"radiusSf": 472.92,
		"data": "ef86ad07aef"
	}`)

	resp := SendRequest(parameters, "POST", "/v1/sigfox/messages")
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestCreateLocation(t *testing.T) {
	parameters := []byte(`
	{
		"sigfoxId": "36DC07",
		"timestamp": 1522152070,
		"latitude": 45.7621635049846,
		"longitude": 4.819336257006938,
		"radius": 4083,
		"spotIt": true
	}`)

	resp := SendRequest(parameters, "POST", "/v1/sigfox/locations")
	assert.Equal(t, http.StatusCreated, resp.Code)
}
