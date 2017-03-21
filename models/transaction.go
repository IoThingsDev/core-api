package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Transaction struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty" valid:"-"`
	UserId    bson.ObjectId `json:"userId" bson:"userId,omitempty" valid:"required"`
	Amount    uint64        `json:"amount" bson:"amount" valid:"required"`
	Date      time.Time     `json:"date" bson:"date" valid:"-"`
	Status    bool          `json:"status" bson:"status" valid:"-"`
	Error     string        `json:"error" bson:"error" valid:"-"`
	CardToken string        `json:"cardToken" bson:"cardToken" valid:"-"`
}

const TransactionsCollection = "transactions"
