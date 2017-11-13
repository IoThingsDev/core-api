package controllers

import (
	"net/http"

	"github.com/adrien3d/things-api/helpers"
	"github.com/adrien3d/things-api/helpers/params"
	"github.com/adrien3d/things-api/models"
	"github.com/adrien3d/things-api/store"
	"gopkg.in/gin-gonic/gin.v1"
)

type DeviceController struct {
}

func NewDeviceController() DeviceController {
	return DeviceController{}
}

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

func (dc DeviceController) GetDevices(c *gin.Context) {
	devices, err := store.GetDevices(c)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, devices)
}

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

func (dc DeviceController) DeleteDevice(c *gin.Context) {
	err := store.DeleteDevice(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (dc DeviceController) GetDevice(c *gin.Context) {
	device, err := store.GetDevice(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, device)
}

func (dc DeviceController) GetLastMessages(c *gin.Context) {
	sigfoxMessages, err := store.GetLastMessages(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, sigfoxMessages)
}

func (dc DeviceController) GetLastLocations(c *gin.Context) {
	locations, err := store.GetLastLocations(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, locations)
}
