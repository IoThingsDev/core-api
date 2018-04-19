package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

/*
	This model is here to implement customization around data format, and data decoding
*/

type SigfoxSyntax struct {
	Id      string `json:"id" bson:"_id,omitempty" valid:"-"`
	FleetId string `json:"fleetId" bson:"fleetId" valid:"-"`
	LastAcc int64  `json:"lastAcc" bson:"lastAcc" valid:"-"`
}

func (d *SigfoxSyntax) BeforeCreate() {
	d.Id = bson.NewObjectId().Hex()
	d.LastAcc = time.Now().Unix()
}

const SigfoxSyntaxesCollection = "sigfoxSyntaxes"
