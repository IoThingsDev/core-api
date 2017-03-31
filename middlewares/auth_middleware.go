package middlewares

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dernise/base-api/models"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func AuthMiddleware(database *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenReader := c.Request.Header.Get("Authorization")

		authHeaderParts := strings.Split(tokenReader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			c.AbortWithError(http.StatusBadRequest, errors.New("Authorization header format must be Bearer {token}"))
			return
		}

		publicKeyFile, _ := ioutil.ReadFile(os.Getenv("BASEAPI_RSA_PUBLIC"))
		publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)

		token, err := jwt.Parse(authHeaderParts[1], func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Error parsing token"))
			return
		}

		if token.Header["alg"] != jwt.SigningMethodRS256.Alg() {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Signing method not valid"))
			return
		}

		if !token.Valid {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Token invalid"))
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)

		session := database.Session.Copy()
		defer session.Close()
		users := database.C(models.UsersCollection).With(session)

		user := models.User{}
		err = users.FindId(bson.ObjectIdHex(claims["id"].(string))).One(&user)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("User not found"))
			return
		}

		c.Set("currentUser", user)
	}
}
