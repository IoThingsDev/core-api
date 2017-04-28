package mongodb

import "gopkg.in/mgo.v2"

type mongo struct {
	*mgo.Database
}

func New(database *mgo.Database) *mongo {
	return &mongo{database}
}
