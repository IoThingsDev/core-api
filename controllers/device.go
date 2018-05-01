package controllers

import (
	"net/http"

	"github.com/IoThingsDev/api/helpers"
	"github.com/IoThingsDev/api/helpers/params"
	"github.com/IoThingsDev/api/models"
	"github.com/IoThingsDev/api/store"
	"gopkg.in/gin-gonic/gin.v1"
)

// Controller that gather all Device related methods
type DeviceController struct {
}

// Initiate a controller for router
func NewDeviceController() DeviceController {
	return DeviceController{}
}

// Create a device
func (dc DeviceController) CreateDevice(c *gin.Context) {
	device := &models.Device{}

	err := c.BindJSON(device)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	if err := store.CreateDevice(c, device); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, device)
}

// Update all devices
func (dc DeviceController) GetDevices(c *gin.Context) {
	devices, err := store.GetDevices(c)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, devices)
}

// Update a specific device
func (dc DeviceController) UpdateDevice(c *gin.Context) {
	device := models.Device{}

	err := c.BindJSON(&device)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	user := store.Current(c)

	changes := params.M{"$set": params.M{"name": device.Name, "userId": user.Id, "bleMac": device.BLEMac, "lastAcc": device.LastAcc, "active": device.Active}}
	err = store.UpdateDevice(
		c,
		c.Param("id"),
		changes,
	)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

// Delete a specific device
func (dc DeviceController) DeleteDevice(c *gin.Context) {
	err := store.DeleteDevice(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

// Getting details from a specific device
func (dc DeviceController) GetDevice(c *gin.Context) {
	device, err := store.GetDevice(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, device)
}

// Getting last messages from a specific device
func (dc DeviceController) GetLastDeviceMessages(c *gin.Context) {
	sigfoxMessages, err := store.GetLastDeviceMessages(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, sigfoxMessages)
}

// Getting last locations from a specific device
func (dc DeviceController) GetLastDeviceLocations(c *gin.Context) {
	locations, err := store.GetLastDeviceLocations(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, locations)
}

// Getting all messages from a specific device
func (dc DeviceController) GetAllDeviceMessages(c *gin.Context) {
	sigfoxMessages, err := store.GetAllDeviceMessages(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, sigfoxMessages)
}

// Getting all locations from a specific device
func (dc DeviceController) GetAllDeviceLocations(c *gin.Context) {
	locations, err := store.GetAllDeviceLocations(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, locations)
}
