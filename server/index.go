package server

import (
	"github.com/dernise/base-api/models"
	"gopkg.in/mgo.v2"
)

func (a API) SetupIndexes() {
	database := a.Database

	// Creates a list of indexes to ensure
	indexes := make(map[*mgo.Collection][]mgo.Index)

	// User indexes
	users := database.C(models.UsersCollection)
	indexes[users] = []mgo.Index{
		{
			Key:    []string{"email"},
			Unique: true,
		},
	}

	for collection, indexes := range indexes {
		for _, index := range indexes {
			err := collection.EnsureIndex(index)

			if err != nil {
				panic("Could not create the indexes")
				return
			}
		}
	}

}
