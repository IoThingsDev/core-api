package controllers

import (
	"net/http"

	"github.com/adrien3d/things-api/store"
	"gopkg.in/gin-gonic/gin.v1"
)

type LocationController struct {
}

func NewLocationController() LocationController {
	return LocationController{}
}

func (lc LocationController) GetAllLocations(c *gin.Context) {
	locations, err := store.GetAllLocations(c)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, locations)
}
