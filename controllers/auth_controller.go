package controllers

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/dernise/pushpal-api/models"
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
	privateKeyFile, _ := ioutil.ReadFile("pushpal.rsa")
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)

	session := ac.mgo.Session.Copy()
	defer session.Close()
	users := ac.mgo.C(models.UsersCollection).With(session)

	userInput := models.User{}
	c.Bind(&userInput) // TODO: HANDLE ERROR

	user := models.User{}
	err := users.Find(bson.M{"username": userInput.Username}).One(&user)
	if err != nil {
		c.Error(err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
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
	token.Claims = claims
	tokenString, err := token.SignedString(privateKey)

	c.JSON(http.StatusOK, gin.H{"status":"success", "data":gin.H{"token":tokenString}})
	return
}
