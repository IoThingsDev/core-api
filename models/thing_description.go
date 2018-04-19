package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

/*
	This model is here to implement thing description exposing
*/

type ThingDescription struct {
	Id       string `json:"id" bson:"_id,omitempty" valid:"-"`
	UserId   string `json:"userId" bson:"userId" valid:"-"`
	Name     string `json:"name" bson:"name" valid:"-"`
	SigfoxId string `json:"sigfoxId" bson:"sigfoxId" valid:"-"`
	BLEMac   string `json:"bleMac" bson:"bleMac" valid:"-"`
	LastAcc  int64  `json:"lastAcc" bson:"lastAcc" valid:"-"`
	Active   bool   `json:"active" bson:"active" valid:"-"`
}

func (d *ThingDescription) BeforeCreate() {
	d.Id = bson.NewObjectId().Hex()
	d.LastAcc = time.Now().Unix()
}

const ThingsDescriptionsCollection = "thingsDescriptions"
