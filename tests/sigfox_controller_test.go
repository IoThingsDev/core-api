package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMessage(t *testing.T) {
	parameters := []byte(`
{
    "sigfoxId": "AABBCC",
    "time": 1523521574,
    "data": "4e2d2f114504301e00411757",
    "mesType":3,
    "snr": 0.00,    
    "avgSnr":0.00,   
    "rssi":0.00,
    "frameNumber": 0,
    "station": "0000",  
    "snr":0.00
}`)

	resp := SendRequest(parameters, "POST", "/v1/sigfox/messages")
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestCreateLocation(t *testing.T) {
	parameters := []byte(`
{
               	"sigfoxId": "17CA54",
               	"spotIt": true,
                "seqNumber": 358, 
                "timestamp": 1523019829,
                "latitude": 45.4540056,
                "longitude": 4.10510950,
                "radius": 3813
}`)

	resp := SendRequest(parameters, "POST", "/v1/sigfox/locations")
	assert.Equal(t, http.StatusCreated, resp.Code)
}
