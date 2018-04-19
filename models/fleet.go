package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

/*
	This model is here to implement thing description exposing
*/

type Fleet struct {
	Id      string   `json:"id" bson:"_id,omitempty" valid:"-"`
	UserId  string   `json:"userId" bson:"userId" valid:"-"`
	Name    string   `json:"name" bson:"name" valid:"-"`
	Devices []Device `json:"devices" bson:"devices" valid:"-"`
	LastAcc int64    `json:"lastAcc" bson:"lastAcc" valid:"-"`
	Active  bool     `json:"active" bson:"active" valid:"-"`
}

func (d *Fleet) BeforeCreate() {
	d.Id = bson.NewObjectId().Hex()
	d.LastAcc = time.Now().Unix()
}

const FleetsCollection = "fleets"
