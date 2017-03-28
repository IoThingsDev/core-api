package server

import (
	"github.com/dernise/base-api/models"
	"gopkg.in/mgo.v2"
)

func (a API) SetupIndexes() {
	database := a.Database

	// Creates a list of indexes to ensure
	collectionIndexes := make(map[*mgo.Collection][]mgo.Index)

	// User indexes
	users := database.C(models.UsersCollection)
	collectionIndexes[users] = []mgo.Index{
		{
			Key:    []string{"email"},
			Unique: true,
		},
	}

	// Transaction indexes
	transactions := database.C(models.TransactionsCollection)
	collectionIndexes[transactions] = []mgo.Index{
		{
			Key: []string{"userId"},
		},
	}

	for collection, indexes := range collectionIndexes {
		for _, index := range indexes {
			err := collection.EnsureIndex(index)

			if err != nil {
				panic(err)
			}
		}
	}

}
