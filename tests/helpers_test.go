package tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"time"

	"github.com/adrien3d/things-api/models"
	"github.com/adrien3d/things-api/server"
	"github.com/adrien3d/things-api/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2/bson"
)

func SendRequest(parameters []byte, method string, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(parameters))
	req.Header.Add("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	return resp
}

func SendRequestWithToken(parameters []byte, method string, url string, authToken string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(parameters))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	return resp
}

func CreateUserAndGenerateToken() (*models.User, string) {
	users := api.Database.C(models.UsersCollection)

	userToken := models.LoginToken{
		Id:         bson.NewObjectId().String(),
		Ip:         "127.0.0.1",
		CreatedAt:  time.Now().Unix(),
		LastAccess: time.Now().Unix(),
	}

	user := models.User{
		Id:        bson.NewObjectId().Hex(),
		Email:     "jeanmichel.lecul@gmail.com",
		Firstname: "Jean-Michel",
		Lastname:  "Lecul",
		Password:  "strongPassword",
		Active:    true,
		StripeId:  "cus_AKlEqL9MjNICJx",
		Admin:     true,
		Tokens: []models.LoginToken{
			userToken,
		},
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	users.Insert(user)

	privateKeyFile, _ := ioutil.ReadFile(api.Config.GetString("rsa_private"))
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)

	token := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodRS256.Alg()))

	claims := make(jwt.MapClaims)
	// TODO: ADD EXPIRATION
	//claims["exp"] = time.Now().Add(time.Hour * time.Duration(settings.Get().JWTExpirationDelta)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["id"] = user.Id

	claims["token"] = userToken.Id

	token.Claims = claims

	tokenString, _ := token.SignedString(privateKey)

	return &user, tokenString
}

func ResetDatabase() {
	api.Database.DropDatabase()
	user, authToken = CreateUserAndGenerateToken()
}

func SetupApi() *server.API {
	api := &server.API{Router: gin.Default(), Config: viper.New()}

	err := api.SetupViper()
	if err != nil {
		panic(err)
	}

	_, err = api.SetupDatabase()
	if err != nil {
		panic(err)
	}

	api.Database.DropDatabase()

	err = api.SetupIndexes()
	if err != nil {
		panic(err)
	}

	// services.SetStripeKeyAndBackend(api.Config)

	api.SetupRedis()

	api.EmailSender = &services.FakeEmailSender{}
	api.SetupRouter()

	return api
}
