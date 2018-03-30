package controllers

import (
	"googlemaps.github.io/maps"
	"net/http"

	"context"
	"fmt"
	"github.com/adrien3d/things-api/helpers"
	"github.com/adrien3d/things-api/models"
	"github.com/adrien3d/things-api/store"
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"strconv"
)

type SigfoxController struct {
}

func NewSigfoxController() SigfoxController {
	return SigfoxController{}
}

func getWifiPosition(ssids string) models.Location {
	fmt.Print("WiFi frame: \t\t\t")
	var wifiLoc models.Location

	ssid1 := ""
	for i := 0; i <= 10; i += 2 {
		if i == 10 {
			ssid1 += fmt.Sprint(string(ssids[i : i+2]))
		} else {
			ssid1 += fmt.Sprint(string(ssids[i:i+2]), ":")
		}
	}
	ssid2 := ""
	for i := 12; i <= 22; i += 2 {
		if i == 22 {
			ssid2 += fmt.Sprint(string(ssids[i : i+2]))
		} else {
			ssid2 += fmt.Sprint(string(ssids[i:i+2]), ":")
		}
	}

	fmt.Print("SSID1: ", ssid1, "\t SSID2:", ssid2, "\t\t\t")
	// TODO: Put Google API Key in config file, like: config.GetString(c, "google_api_key")
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyCN0Z78M1sIT6c2H8PL0KaaFmjkBUE4avQ"))
	if err != nil {
		log.Fatalf("API connection fatal error: %s", err)
	}
	r := &maps.GeolocationRequest{
		ConsiderIP: false,
		WiFiAccessPoints: []maps.WiFiAccessPoint{{
			MACAddress: ssid1,
		}, {
			MACAddress: ssid2,
		}},
	}
	resp, err := c.Geolocate(context.Background(), r)
	if err != nil {
		log.Println("Google Maps Geolocation Request error: %s", err)
	}

	//fmt.Println(resp)

	wifiLoc.Latitude = resp.Location.Lat
	wifiLoc.Longitude = resp.Location.Lng
	wifiLoc.Radius = resp.Accuracy
	wifiLoc.SpotIt = false
	wifiLoc.GPS = false
	wifiLoc.WiFi = true
	fmt.Println(wifiLoc)
	return wifiLoc
}

func decodeGPSFrame(msg models.SigfoxMessage) (models.Location, float64) {
	fmt.Print("GPS frame: \t\t\t")
	var gpsLoc models.Location
	var temperature float64
	var latitude, longitude float64
	var latDeg, latMin, latSec float64
	var lngDeg, lngMin, lngSec float64

	isNorth, isEast := false, false
	if string(msg.Data[0:2]) == "4e" {
		isNorth = true
	}
	if string(msg.Data[10:12]) == "45" {
		isEast = true
	}

	if isNorth {
		fmt.Print("N:")
	} else {
		fmt.Print("S:")
	}

	valLatDeg, _ := strconv.ParseInt(msg.Data[2:4], 16, 8)
	latDeg = float64(valLatDeg)
	valLatMin, _ := strconv.ParseInt(msg.Data[4:6], 16, 8)
	latMin = float64(valLatMin)
	valLatSec, _ := strconv.ParseInt(msg.Data[6:8], 16, 8)
	latSec = float64(valLatSec)
	fmt.Print(latDeg, "° ", latMin, "m ", latSec, "s\t")

	latitude = float64(latDeg) + float64(latMin/60) + float64(latSec/3600)

	if isEast {
		fmt.Print("E:")
	} else {
		fmt.Print("W:")
	}
	valLngDeg, _ := strconv.ParseInt(msg.Data[10:12], 16, 8)
	lngDeg = float64(valLngDeg)
	valLngMin, _ := strconv.ParseInt(msg.Data[12:14], 16, 8)
	lngMin = float64(valLngMin)
	valLngSec, _ := strconv.ParseInt(msg.Data[14:16], 16, 8)
	lngSec = float64(valLngSec)
	fmt.Print(lngDeg, "° ", lngMin, "m ", lngSec, "s")

	longitude = float64(lngDeg) + float64(lngMin/60) + float64(lngSec/3600)

	fmt.Print("\t\t\t Lat: ", latitude, "\t Lng:", longitude)
	// Populating returned location
	gpsLoc.Latitude = latitude
	gpsLoc.Longitude = longitude
	gpsLoc.SpotIt = false
	gpsLoc.GPS = true
	gpsLoc.WiFi = false

	temperature, err := strconv.ParseFloat(msg.Data[20:22], 64)
	if err != nil {
		fmt.Println("Error while converting temperature main")
	}
	dec, err := strconv.ParseFloat(msg.Data[22:24], 64)
	if err != nil {
		fmt.Println("Error while converting temperature decimal")
	}

	temperature += dec * 0.01

	fmt.Println("\t\t", gpsLoc, "\t", temperature)
	return gpsLoc, temperature
}

func (sc SigfoxController) CreateMessage(c *gin.Context) {
	sigfoxMessage := &models.SigfoxMessage{}

	err := c.BindJSON(sigfoxMessage)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	if sigfoxMessage.MesType == 3 {
		computedLocation := &models.Location{}
		//var computedLocation models.Location

		if (string(sigfoxMessage.Data[0:2]) == "4e") || (string(sigfoxMessage.Data[0:2]) == "53") {
			decodedGPSFrame, decodedTemperature := decodeGPSFrame(*sigfoxMessage)
			sigfoxMessage.Data4 = decodedTemperature
			computedLocation = &decodedGPSFrame
			fmt.Println("Wisol GPS Frame, contaning: ", computedLocation)
		} else {
			decodedWifiFrame := getWifiPosition(sigfoxMessage.Data)
			computedLocation = &decodedWifiFrame
			fmt.Println("Wisol WiFi Frame, contaning: ", computedLocation)
			//store.CreateLocation(context.Background(), &wifiLoc)
		}

		computedLocation.SigfoxId = sigfoxMessage.SigfoxId

		//err = store.CreateLocationWithMessage(c, computedLocation, sigfoxMessage)
		err = store.CreateLocation(c, computedLocation)
		fmt.Println("Computed location created")
		if err != nil {
			fmt.Println("Error while creating computed location")
			c.Error(err)
			c.Abort()
			return
		}

		sigfoxMessage.Data1 = computedLocation.Latitude
		sigfoxMessage.Data2 = computedLocation.Longitude
		sigfoxMessage.Data3 = computedLocation.Radius

		if computedLocation.SpotIt {
			sigfoxMessage.Data4 = 1
		} else if computedLocation.GPS {
			sigfoxMessage.Data5 = 1
		} else if computedLocation.WiFi {
			sigfoxMessage.Data6 = 1
		}

		err = store.CreateMessage(c, sigfoxMessage)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}
	c.JSON(http.StatusCreated, sigfoxMessage)
}
