package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"gopkg.in/mgo.v2"
)

type UserController struct {
	c *mgo.Collection
}


func NewUserController(c *mgo.Collection) *UserController {
	return &UserController{
		c,
	}
}

func (uc UserController) GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "data":c.Param("id")})

}

func (uc UserController) CreateUser(c *gin.Context) {

}