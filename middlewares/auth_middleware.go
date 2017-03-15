package middlewares

import (
	"gopkg.in/gin-gonic/gin.v1"
	"strings"
	"net/http"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenReader := c.Request.Header.Get("Authorization")

		authHeaderParts := strings.Split(tokenReader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			c.AbortWithError(http.StatusBadRequest, errors.New("Authorization header format must be Bearer {token}"))
			return
		}

		publicKeyFile, err := ioutil.ReadFile("pushpal.rsa.pub")
		publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)

		token, err := jwt.Parse(authHeaderParts[1], func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Error parsing token"))
			return
		}

		if token.Header["alg"] != jwt.SigningMethodRS256.Alg() {
			log.Println(jwt.SigningMethodRS256)
			c.AbortWithError(http.StatusUnauthorized, errors.New("Signing method not valid"))
			return
		}

		if !token.Valid {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Token invalid"))
			return
		}
	}
}
