package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"gopkg.in/mgo.v2"
	"github.com/dernise/pushpal-api/models"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/asaskevich/govalidator"
	"errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"bytes"
	"github.com/sendgrid/rest"
	"text/template"
	"github.com/spf13/viper"
	"github.com/dernise/pushpal-api/helpers"
)

type UserController struct {
	mgo    *mgo.Database
	config *viper.Viper
}

func NewUserController(mgo *mgo.Database, config *viper.Viper) *UserController {
	return &UserController{
		mgo,
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
	from := mail.NewEmail("Pushpal", "no-reply@pushpal.io")
	subject := "Welcome to Pushpal! Confirm your email"
	to := mail.NewEmail(user.Firstname, user.Email)


	url := "Please confirm your email address by clicking on the following link: http://localhost:4000/users/{{.Id.Hex}}/activate/{{.ActivationKey}}"
	buffer := new(bytes.Buffer)
	template := template.Must(template.New("emailTemplate").Parse(url))
	template.Execute(buffer, user)

	content := mail.NewContent("text/plain", buffer.String())
	m := mail.NewV3MailInit(from, subject, to, content)
	request := sendgrid.GetRequest(uc.config.GetString("sendgrid.apiKey"), "/v3/mail/send", "")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	return response, err
}
