package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

type (
	UserController struct{}
)

func NewUserController() *UserController {
	return &UserController{}
}

func (uc UserController) GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "data":c.Param("id")})
}

func (uc UserController) CreateUser(c *gin.Context) {

}