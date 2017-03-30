package controllers

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AuthController struct {
	mgo    *mgo.Database
	config *viper.Viper
}

func NewAuthController(mgo *mgo.Database, config *viper.Viper) *AuthController {
	return &AuthController{
		mgo,
		config,
	}
}

func (ac AuthController) Authentication(c *gin.Context) {
	privateKeyFile, _ := ioutil.ReadFile(ac.config.GetString("rsa_private"))
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)

	session := ac.mgo.Session.Copy()
	defer session.Close()
	users := ac.mgo.C(models.UsersCollection).With(session)

	userInput := models.User{}
	c.Bind(&userInput) // TODO: HANDLE ERROR

	user := models.User{}
	err := users.Find(bson.M{"email": userInput.Email}).One(&user)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "User does not exist"))
		return
	}

	if !user.Active {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_needs_activation", "User needs to be activated"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("incorrect_password", "Password is not correct"))
		return
	}

	token := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodRS256.Alg()))
	claims := make(jwt.MapClaims)
	// TODO: ADD EXPIRATION
	//claims["exp"] = time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["id"] = user.Id
	token.Claims = claims
	tokenString, err := token.SignedString(privateKey)

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "users": user.Sanitize()})
}
