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
)

type UserController struct {
	mgo *mgo.Database
}


func NewUserController(mgo *mgo.Database) *UserController {
	return &UserController{
		mgo,
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

	c.JSON(http.StatusOK, gin.H{"status": "success", "data":user})
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

	err = users.Insert(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Failed to insert the user"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data":user})
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

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data":list})
}
