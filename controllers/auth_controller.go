package controllers

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/pushpal-api/models"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"time"
)

type AuthController struct {
	mgo *mgo.Database
}

func NewAuthController(mgo *mgo.Database) *AuthController {
	return &AuthController{
		mgo,
	}
}


func (ac AuthController) Authentication(c *gin.Context) {
	privateKey, _ := ioutil.ReadFile("pushpal.rsa")

	session := ac.mgo.Session.Copy()
	defer session.Close()
	users := ac.mgo.C(models.UsersCollection).With(session)

	password := []byte(c.Param("password"))
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		c.Error(err)
		return
	}

	user := models.User{}
	err = users.Find(bson.M{"username": c.Param("username"), "password": string(hashedPassword)}).One(&user)
	if err != nil {
		c.Error(err)
		return
	}

	token := jwt.New(jwt.GetSigningMethod("RS256"))

	claims := make(jwt.MapClaims)
	// TODO: ADD EXPIRATION
	//claims["exp"] = time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["id"] = user.Id
	claims["username"] = user.Username
	token.Claims = claims

	tokenString, _ := token.SignedString(privateKey)

	c.JSON(http.StatusOK, gin.H{"status":"success", "data":gin.H{"token":tokenString}})
	return
}
