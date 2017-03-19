package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"gopkg.in/mgo.v2"
	"github.com/dernise/base-api/models"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/asaskevich/govalidator"
	"errors"
	"github.com/spf13/viper"
	"github.com/dernise/base-api/helpers"
	"github.com/sendgrid/rest"
	"bytes"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"github.com/dernise/base-api/services"
)


type UserController struct {
	mgo    *mgo.Database
	config *viper.Viper
	emailSender services.EmailSender
}


func NewUserController(mgo *mgo.Database, config *viper.Viper, emailSender services.EmailSender) *UserController {
	return &UserController{
		mgo,
		config,
		emailSender,
	}
}

func (uc UserController) GetUser(c *gin.Context) {
	session := uc.mgo.Session.Copy()
	defer session.Close()
	users := uc.mgo.C(models.UsersCollection).With(session)

	user := models.User{}
	err := users.Find(bson.M{"_id": bson.ObjectIdHex(c.Param("id"))}).One(&user)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, errors.New("User not found"))
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
		c.AbortWithError(http.StatusBadRequest, errors.New("Failed to bind the body data"))
		return
	}

	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	count, _ := users.Find(bson.M{"email": user.Email}).Count()
	if count > 0 {
		c.AbortWithError(http.StatusConflict, errors.New("User already exists"))
		return
	}

	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Failed to generate the encrypted password"))
		return
	}

	user.Active = false
	user.ActivationKey = helpers.RandomString(20)

	user.Id = bson.NewObjectId()

	err = users.Insert(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Failed to insert the user"))
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
		c.AbortWithError(http.StatusNotFound, errors.New("Users not found"))
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
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User has been activated"})
}

func (uc UserController) SendActivationEmail(user *models.User) (*rest.Response, error) {
	subject := "Welcome to base! Confirm your email"
	to := mail.NewEmail(user.Firstname, user.Email)

	url := "Please confirm your email address by clicking on the following link: http://localhost:4000/users/{{.Id.Hex}}/activate/{{.ActivationKey}}"
	buffer := new(bytes.Buffer)
	template := template.Must(template.New("emailTemplate").Parse(url))
	template.Execute(buffer, user)

	response, err := uc.emailSender.SendEmail([]*mail.Email{ to }, "text/plain", subject, buffer.String())
	return response, err
}
