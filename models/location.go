package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Location struct {
	Id        string  `json:"id" bson:"_id,omitempty" valid:"-"`
	SigfoxId  string  `json:"sigfoxId" bson:"sigfoxId" valid:"-"`
	Timestamp int64   `json:"timestamp" bson:"timestamp" valid:"-"`
	Latitude  float64 `json:"latitude" bson:"latitude" valid:"-"`
	Longitude float64 `json:"longitude" bson:"longitude" valid:"-"`
	Radius    float64 `json:"radius" bson:"radius" valid:"-"`
	SpotIt    bool    `json:"spotIt" bson:"spotIt" valid:"-"`
	GPS       bool    `json:"gps" bson:"gps" valid:"-"`
	WiFi      bool    `json:"WiFi" bson:"WiFi" valid:"-"`
}

type LastLocation struct {
	DeviceId   string   `json:"id" bson:"_id,omitempty" valid:"-"`
	DeviceName string   `json:"name" bson:"name"`
	Location   Location `json:"location" bson:"location"`
}

func (l *Location) BeforeCreate() {
	l.Id = bson.NewObjectId().Hex()
	l.Timestamp = time.Now().Unix()
}

const LocationsCollection = "locations"
