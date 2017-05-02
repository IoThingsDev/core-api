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

type UserController struct {
	emailSender services.EmailSender
}

func NewUserController(emailSender services.EmailSender) UserController {
	return UserController{
		emailSender,
	}
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
	uc.sendActivationEmail(appName, user)

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

// Checks for a user that matches an email, and sends a reset mail
/*func (uc UserController) ResetPasswordRequest(c *gin.Context) {
	session := uc.mgo.Session.Copy()
	defer session.Close()
	users := uc.mgo.C(models.UsersCollection).With(session)

	err := uc.redis.UpdateEmailRateLimit(c.ClientIP())
	if err != nil {
		c.AbortWithError(http.StatusTooManyRequests, helpers.ErrorWithCode("too_many_requests", "You sent too many requests on this endpoint."))
		return

	}

	user := models.User{}
	err = c.Bind(&user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	err = users.Find(bson.M{"email": user.Email}).One(&user)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "User not found"))
		return
	}

	resetKey := helpers.RandomString(20)

	err = users.UpdateId(user.Id, bson.M{"$set": bson.M{"resetKey": resetKey}})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("update_resetKey_failed", "Failed to update the resetKey"))
		return
	}

	uc.sendResetPasswordRequestEmail(&user)
}

// TODO: GET idUser, resetToken : FormResetPassword : renderForm to post to ResetPassword

// POST userid, resetKey and new password to change them
func (uc UserController) ResetPassword(c *gin.Context) {
	session := uc.mgo.Session.Copy()
	defer session.Close()
	users := uc.mgo.C(models.UsersCollection).With(session)

	userId := c.Param("id")
	resetKey := c.PostForm("resetKey")
	newPassword := c.PostForm("newPassword")

	user := models.User{}

	password := []byte(newPassword)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	hashedPasswordStore := string(hashedPassword)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("encryption_failed", "Failed to generate the encrypted password"))
		return
	}

	_, err = users.Find(bson.M{"$and": []bson.M{{"_id": bson.ObjectIdHex(userId)}, {"resetKey": resetKey}}}).Apply(mgo.Change{Update: bson.M{"$set": bson.M{"password": hashedPasswordStore}}}, &user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("update_password_failed", "Failed to update the password"))
		return
	}

	uc.sendResetPasswordDoneEmail(&user)
}*/

func (uc UserController) sendActivationEmail(appName string, user *models.User) {
	subject := "Welcome to " + appName + "! Account confirmation"
	templateLink := "./templates/mail_activate_account.html"
	uc.emailSender.SendEmailFromTemplate(user, subject, templateLink)
}

/*func (uc UserController) sendResetPasswordRequestEmail(user *models.User) {
	appName := uc.config.GetString("sendgrid_name")
	subject := "Reset your " + appName + " Account password"
	templateLink := "./templates/mail_reset_password_request.html"
	uc.emailSender.SendEmailFromTemplate(user, subject, templateLink)
}

func (uc UserController) sendResetPasswordDoneEmail(user *models.User) {
	appName := uc.config.GetString("sendgrid_name")
	subject := "Your " + appName + " password has been reset"
	templateLink := "./templates/mail_reset_password_done.html"
	uc.emailSender.SendEmailFromTemplate(user, subject, templateLink)
}*/
