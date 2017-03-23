package controllers

import (
	"github.com/asaskevich/govalidator"
	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
	"github.com/dernise/base-api/services"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type UserController struct {
	mgo         *mgo.Database
	emailSender services.EmailSender
	config      *viper.Viper
}

func NewUserController(mgo *mgo.Database, emailSender services.EmailSender, config *viper.Viper) *UserController {
	return &UserController{
		mgo,
		emailSender,
		config,
	}
}
func (uc UserController) GetUser(c *gin.Context) {
	session := uc.mgo.Session.Copy()
	defer session.Close()
	users := uc.mgo.C(models.UsersCollection).With(session)

	user := models.User{}
	err := users.FindId(bson.ObjectIdHex(c.Param("id"))).One(&user)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "User not found"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
}


func (uc UserController) GetUsers(c *gin.Context) {
	session := uc.mgo.Session.Copy()
	defer session.Close()
	users := uc.mgo.C(models.UsersCollection).With(session)

	list := []models.User{}
	err := users.Find(nil).All(&list)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "Users not found"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": list})
}

func (uc UserController) CreateUser(c *gin.Context) {
	session := uc.mgo.Session.Copy()
	defer session.Close()
	users := uc.mgo.C(models.UsersCollection).With(session)

	user := models.User{}
	err := c.Bind(&user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	count, _ := users.Find(bson.M{"email": user.Email}).Count()
	if count > 0 {
		c.AbortWithError(http.StatusConflict, helpers.ErrorWithCode("user_already_exists", "User already exists"))
		return
	}

	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("encryption_failed", "Failed to generate the encrypted password"))
		return
	}

	user.Active = false
	user.ActivationKey = helpers.RandomString(20)
	user.ResetKey = helpers.RandomString(20)

	user.StripeId = ""

	user.Id = bson.NewObjectId()

	uc.SendActivationEmail(&user)

	err = users.Insert(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("creation_failed", "Failed to insert the user"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "User created"})
}

func (uc UserController) ActivateUser(c *gin.Context) {
	session := uc.mgo.Session.Copy()
	defer session.Close()
	users := uc.mgo.C(models.UsersCollection).With(session)

	userId := c.Param("id")
	activationKey := c.Param("activationKey")

	err := users.Update(bson.M{"$and": []bson.M{{"_id": bson.ObjectIdHex(userId)}, {"activationKey": activationKey}}}, bson.M{"$set": bson.M{"active": true}})
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("activation_failed", "Couldn't find the user to activate"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User has been activated"})
}

//Checks for a user that matches an email, and send a reset mail
func (uc UserController) ResetPasswordRequest(c *gin.Context) {
	session := uc.mgo.Session.Copy()
	defer session.Close()
	users := uc.mgo.C(models.UsersCollection).With(session)


	user := models.User{}
	err := c.Bind(&user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data"))
		return
	}

	err = users.Find(bson.M{"email": user.Email}).One(&user)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "User not found"))
		return
	}

	uc.SendResetPasswordRequestEmail(&user)
}

//GET userid & resetKey to be able to change password
func (uc UserController) ControlResetKey(c *gin.Context) {
	session := uc.mgo.Session.Copy()
	defer session.Close()
	users := uc.mgo.C(models.UsersCollection).With(session)

	userId := c.Param("id")
	resetKey := c.Param("resetKey")

	err := users.Update(bson.M{"$and": []bson.M{{"_id": bson.ObjectIdHex(userId)}, {"resetKey": resetKey}}}, bson.M{"$set": bson.M{"active": true}})
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("activation_failed", "Couldn't find the user to reset"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User can change his password"})
}

//POST userid, resetKey and new password to change them
func (uc UserController) ResetPassword(c *gin.Context) {

}

func (uc UserController) SendActivationEmail(user *models.User) {
	appName := uc.config.GetString("sendgrid_name")
	subject := "Welcome to " + appName + "! Account confirmation"
	templateLink := "./templates/mail_activate_account.html"
	uc.emailSender.SendEmailFromTemplate(user, subject, templateLink)
}

func (uc UserController) SendResetPasswordRequestEmail(user *models.User) {
	appName := uc.config.GetString("sendgrid_name")
	subject := "Reset your " + appName + " Account password"
	templateLink := "./templates/mail_reset_password_request.html"
	uc.emailSender.SendEmailFromTemplate(user, subject, templateLink)
}

func (uc UserController) SendResetPasswordDoneEmail(user *models.User) {
	appName := uc.config.GetString("sendgrid_name")
	subject := "Your " + appName + " has been reset"
	templateLink := "./templates/mail_reset_password_done.html"
	uc.emailSender.SendEmailFromTemplate(user, subject, templateLink)
}
