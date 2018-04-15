package middlewares

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adrien3d/things-api/helpers"
	"github.com/adrien3d/things-api/helpers/params"
	"github.com/adrien3d/things-api/models"
	"github.com/adrien3d/things-api/services"
	"github.com/adrien3d/things-api/store"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gin-gonic/gin.v1"
)

func AuthMiddleware() gin.HandlerFunc {
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
			c.AbortWithError(http.StatusUnauthorized, errors.New("LoginToken invalid"))
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)

		user := &models.User{}

		// Gets the user from the redis store
		hasFetchedRedis := true
		err = services.GetRedis(c).GetValueForKey(claims["id"].(string), &user)
		if err != nil {
			hasFetchedRedis = false
			user, err = store.FindUserById(c, claims["id"].(string))
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("token_not_valid", "This token isn't valid"))
				return
			}
			err = services.GetRedis(c).SetValueForKey(user.Id, &user)

		}

		// Check if the token is still valid in the database
		loginToken := claims["token"].(string)
		tokenIndex, hasToken := user.HasToken(loginToken)
		if !hasToken {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("token_invalidated", "This token isn't valid anymore"))
			return
		}

		c.Set(store.CurrentKey, user)
		c.Set(store.LoginTokenKey, loginToken)

		if !hasFetchedRedis {
			store.UpdateUser(c, params.M{"$set": params.M{"tokens." + strconv.Itoa(tokenIndex) + ".last_access": time.Now().Unix()}})
		}

		c.Next()
	}
}
