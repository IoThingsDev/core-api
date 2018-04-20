package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Fleet struct {
	Id        string   `json:"id" bson:"_id,omitempty" valid:"-"`
	UserId    string   `json:"userId" bson:"userId" valid:"-"`
	Name      string   `json:"name" bson:"name"`
	DeviceIds []string `json:"device_ids" bson:"device_ids"`
	LastAcc   int64    `json:"lastAcc" bson:"lastAcc" valid:"-"`
	Active    bool     `json:"active" bson:"active" valid:"-"`
}

func (d *Fleet) BeforeCreate() {
	d.Id = bson.NewObjectId().Hex()
	d.LastAcc = time.Now().Unix()
}

const FleetsCollection = "fleets"
