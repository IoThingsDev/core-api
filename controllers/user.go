package controllers

import (
	"net/http"

	"github.com/IoThingsDev/api/config"
	"github.com/IoThingsDev/api/helpers"
	"github.com/IoThingsDev/api/models"
	"github.com/IoThingsDev/api/services"
	"github.com/IoThingsDev/api/store"

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

	c.JSON(http.StatusOK, user.Sanitize())
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

	c.JSON(http.StatusCreated, user.Sanitize())
}

func (uc UserController) ActivateUser(c *gin.Context) {
	if err := store.ActivateUser(c, c.Param("activationKey"), c.Param("id")); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "user_activated.tmpl", gin.H{
		"url": config.GetString(c, "FRONT_URL"),
	})
	//c.JSON(http.StatusOK, gin.H{"status": "success", "message": "You successfully activated your account."})
}

// TODO: Reset Password: rate limit: 6h, 10 total
// Request, sending mail : If your email address exists in our database, you will receive a password recovery link at your email address in a few minutes.
// Handle mail-click (with token checking), clear current password
// Modify password (2 repeats), can be merged with modify in User management UI
