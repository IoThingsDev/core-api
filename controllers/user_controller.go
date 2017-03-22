package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"gopkg.in/mgo.v2"
	"github.com/dernise/base-api/models"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/asaskevich/govalidator"
	"github.com/dernise/base-api/helpers"
	"bytes"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"github.com/dernise/base-api/services"
	"github.com/sendgrid/rest"
	"io/ioutil"
	"github.com/spf13/viper"
	"errors"
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
	err := users.Find(bson.M{"_id": bson.ObjectIdHex(c.Param("id"))}).One(&user)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "User not found"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
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

	user.Id = bson.NewObjectId()

	err = users.Insert(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("creation_failed", "Failed to insert the user"))
		return
	}

	uc.SendActivationEmail(&user)

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "User created"})
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

func (uc UserController) ActivateUser(c *gin.Context) {
	session := uc.mgo.Session.Copy()
	defer session.Close()
	users := uc.mgo.C(models.UsersCollection).With(session)

	key := c.Param("key")
	userId := c.Param("id")

	err := users.Update(bson.M{"$and": []bson.M{{"_id": bson.ObjectIdHex(userId)}, {"activationKey": key}}}, bson.M{"$set": bson.M{"active": true}})
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("activation_failed", "Couldn't find the user to activate"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User has been activated"})
}

func (uc UserController) SendActivationEmail(user *models.User) (*rest.Response, error) {
	type Data struct {
		User        *models.User
		HostAddress string
		AppName     string
	}

	serverName := uc.config.GetString("sendgrid.address")
	appName := uc.config.GetString("sendgrid.name")
	hostname := uc.config.GetString("host.address")

	subject := "Welcome to " + serverName + "! Account confirmation"

	to := mail.NewEmail(user.Firstname, user.Email)

	buffer := new(bytes.Buffer)

	file, err := ioutil.ReadFile("./templates/mail_confirm_account.html")

	if err != nil || len(string(file)) == 0 {
		return nil, err
	}
	
	htmlTemplate := template.Must(template.New("emailTemplate").Parse(string(file)))
	data := Data{User: user, HostAddress: hostname, AppName: appName}
	htmlTemplate.Execute(buffer, data)

	response, err := uc.emailSender.SendEmail([]*mail.Email{to }, "text/html", subject, buffer.String())

	return response, err
}

//TODO: func (uc UserController) SendResetPasswordEmail(user *models.User) (*rest.Response, error)
//TODO: func (uc UserController) ResetPassword(c *gin.Context)
