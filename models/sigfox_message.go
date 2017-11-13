package models

import (
	"gopkg.in/mgo.v2/bson"
)

type SigfoxMessage struct {
	Id          string  `json:"id" bson:"_id,omitempty" valid:"-"`
	SigfoxId    string  `json:"sigfoxId" bson:"sigfoxId" valid:"-"`
	FrameNumber uint    `json:"frameNumber" bson:"frameNumber" valid:"-"` //Device : (daily frames under 140)
	Timestamp   int64   `json:"timestamp" bson:"timestamp" valid:"-"`     //Sigfox : time
	Station     string  `json:"station" bson:"station" valid:"-"`         //Sigfox : station
	Snr         float64 `json:"snr" bson:"snr" valid:"-"`                 //Sigfox : snr
	AvgSnr      float64 `json:"avgSnr" bson:"avgSnr" valid:"-"`           //Sigfox : avgSnr
	Rssi        float64 `json:"rssi" bson:"rssi" valid:"-"`               //Sigfox : rssi
	Data        string  `json:"data" bson:"data" valid:"-"`               //Sigfox : data
	data1		float64 `json:"data1" bson:"data1" valid:"-"`           //Device : data1
	data2		float64 `json:"data2" bson:"data2" valid:"-"`           //Device : data2
	data3   	float64 `json:"data3" bson:"data3" valid:"-"`    	 //Device : data3
	data4		float64 `json:"data4" bson:"data4" valid:"-"`         //Device : data4
	data5		float64 `json:"data5" bson:"data5" valid:"-"`         //Device : data5
	data6		float64 `json:"data6" bson:"data6" valid:"-"`         //Device : data6
}

func (l *SigfoxMessage) BeforeCreate() {
	l.Id = bson.NewObjectId().Hex()
}

const SigfoxMessagesCollection = "sigfox_messages"
