package controllers

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"gopkg.in/mgo.v2"
	"github.com/dernise/pushpal-api/models"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
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
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data":user})
	return
}

func (uc UserController) CreateUser(c *gin.Context) {
	session := uc.mgo.Session.Copy()
	defer session.Close()
	users := uc.mgo.C(models.UsersCollection).With(session)

	user := models.User{}
	err := c.Bind(&user)
	if err != nil {
		c.Error(err)
		return
	}

	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	if err != nil {
		c.Error(err)
		return
	}

	err = users.Insert(user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data":user})
	return
}

func (uc UserController) GetUsers(c *gin.Context) {
	session := uc.mgo.Session.Copy()
	defer session.Close()
	users := uc.mgo.C(models.UsersCollection).With(session)

	list := []models.User{}
	err := users.Find(nil).All(&list)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data":list})
	return
}
