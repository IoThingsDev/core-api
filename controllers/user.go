package controllers

import (
	"net/http"

	"github.com/dernise/base-api/config"
	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
	"github.com/dernise/base-api/services"
	"github.com/dernise/base-api/store"
	"gopkg.in/gin-gonic/gin.v1"
)

type UserController struct{}

func NewUserController() UserController {
	return UserController{}
}
func (uc UserController) GetUser(c *gin.Context) {
	user, err := store.FindUserById(c, c.Param("id"))

	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": user.Sanitize()})
}

func (uc UserController) CreateUser(c *gin.Context) {
	user := &models.User{}

	if err := c.BindJSON(user); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	if err := store.CreateUser(c, user); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	appName := config.GetString(c, "sendgrid_name")
	subject := "Welcome to " + appName + "! Account confirmation"
	templateLink := "./templates/mail_activate_account.html"
	services.GetEmailSender(c).SendEmailFromTemplate(user, subject, templateLink)

	c.JSON(http.StatusCreated, gin.H{"users": user.Sanitize()})
}

func (uc UserController) ActivateUser(c *gin.Context) {
	if err := store.ActivateUser(c, c.Param("activationKey"), c.Param("id")); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}
