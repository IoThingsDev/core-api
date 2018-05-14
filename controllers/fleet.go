package controllers

import (
	"github.com/IoThingsDev/api/helpers"
	"github.com/IoThingsDev/api/models"
	"github.com/IoThingsDev/api/store"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"os"
)

//FleetController handles things related to fleets
type FleetController struct {
}

//NewFleetController instantiate a FleetController
func NewFleetController() FleetController {
	return FleetController{}
}

//Getting all Temperatures from all devices
func (fc FleetController) GetTemperatures(c *gin.Context) {
	id := c.Param("id")

	sigfoxMessages, err := store.GetLastFleetMessages(c, id)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	var formattedMessages []gin.H
	for _, message := range sigfoxMessages {
		formattedMessage := message.FormatData("temperature")
		formattedMessage["deviceId"] = message.SigfoxId
		formattedMessage["href"] = os.Getenv("BASEAPI_BASE_URL") + "/devices/" + message.SigfoxId + "/description"
		formattedMessages = append(formattedMessages, formattedMessage)
	}

	c.JSON(http.StatusOK, formattedMessages)
}

//Getting mean Temperatures from all devices
func (fc FleetController) GetAverageTemperature(c *gin.Context) {
	id := c.Param("id")

	sigfoxMessages, err := store.GetLastFleetMessages(c, id)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	sum := float64(0)
	for _, message := range sigfoxMessages {
		sum += message.Data2
	}

	average := sum / float64(len(sigfoxMessages))

	c.JSON(http.StatusOK, gin.H{
		"value": average,
	})
}

//Getting all fleets
func (fc FleetController) GetFleets(c *gin.Context) {
	fleets, err := store.GetAllFleets(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, fleets)
}

//Getting a single fleet
func (fc FleetController) GetFleetById(c *gin.Context) {
	id := c.Param("id")

	fleet, err := store.GetFleetById(c, id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, fleet)
}

//Creating a fleet
func (fc FleetController) CreateFleet(c *gin.Context) {
	fleet := &models.Fleet{}

	if err := c.BindJSON(fleet); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	if err := store.CreateFleet(c, fleet); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, fleet)
}

//Delete a fleet
func (fc FleetController) DeleteFleet(c *gin.Context) {
	err := store.DeleteFleet(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}
