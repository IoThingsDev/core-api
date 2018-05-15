package controllers

import (
	"net/http"

	"github.com/IoThingsDev/api/helpers"
	"github.com/IoThingsDev/api/models"
	"github.com/IoThingsDev/api/store"
	"gopkg.in/gin-gonic/gin.v1"
)

type LocationController struct {
}

func NewLocationController() LocationController {
	return LocationController{}
}

// Create a single Location
func (lc LocationController) CreateLocation(c *gin.Context) {
	location := &models.Location{}

	err := c.BindJSON(location)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	err = store.CreateLocation(c, location)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, location)
}

// Getting Last Locations from all devices of a user
func (lc LocationController) GetLastDevicesLocations(c *gin.Context) {
	locations, err := store.GetLastDevicesLocations(c)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, locations)
}
