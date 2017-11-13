package server

import (
	"github.com/adrien3d/things-api/models"

	"gopkg.in/mgo.v2"
)

func (a *API) SetupIndexes() error {
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

	locations := database.C(models.LocationsCollection)
	collectionIndexes[locations] = []mgo.Index{
		{
			Key: []string{"sigfoxId"},
		},
	}

	devices := database.C(models.DevicesCollection)
	collectionIndexes[devices] = []mgo.Index{
		{
			Key: []string{"userId"},
		},
		{
			Key: []string{"tokens._id"},
		},
	}

	for collection, indexes := range collectionIndexes {
		for _, index := range indexes {
			err := collection.EnsureIndex(index)

			if err != nil {
				return err
			}
		}
	}
	return nil
}
