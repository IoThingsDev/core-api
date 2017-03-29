package tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/dernise/base-api/models"
	"github.com/dernise/base-api/server"
	"github.com/dgrijalva/jwt-go"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type FakeEmailSender struct {
	to          []*mail.Email
	contentType string
	subject     string
	body        string
}

func (f FakeEmailSender) SendEmail(to []*mail.Email, contentType, subject, body string) (*rest.Response, error) {
	f.to, f.contentType, f.subject, f.body = to, contentType, subject, body
	return &rest.Response{StatusCode: http.StatusOK, Body: "Everything's fine Jean-Miche", Headers: nil}, nil
}

func (f FakeEmailSender) SendEmailFromTemplate(user *models.User, subject string, templateLink string) (*rest.Response, error) {
	return &rest.Response{StatusCode: http.StatusOK, Body: "Everything's fine Jean-Miche", Headers: nil}, nil
}

func SetupApi() *server.API {
	api := server.API{Router: gin.Default(), Config: viper.New()}

	err := api.SetupViper()
	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial(api.Config.GetString("db_host"))
	if err != nil {
		panic(err)
	}

	api.Database = session.DB(api.Config.GetString("db_name"))
	api.Database.DropDatabase()

	err = api.SetupIndexes()
	if err != nil {
		panic(err)
	}
	api.EmailSender = &FakeEmailSender{}
	api.SetupRouter()

	return &api
}

func SendRequest(api *server.API, parameters []byte, method string, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(parameters))
	req.Header.Add("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	return resp
}

func SendRequestWithToken(api *server.API, parameters []byte, method string, url string, jwtToken string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(parameters))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+jwtToken)
	resp := httptest.NewRecorder()
	api.Router.ServeHTTP(resp, req)
	return resp
}

func CreateUserAndGenerateToken(api *server.API) (*models.User, string) {
	users := api.Database.C(models.UsersCollection)

	user := models.User{
		Id:        bson.NewObjectId(),
		Email:     "jeanmichel.lecul@gmail.com",
		Firstname: "Jean-Michel",
		Lastname:  "Lecul",
		Password:  "strongPassword",
		Active:    true,
		StripeId:  "cus_AKlEqL9MjNICJx",
	}

	users.Insert(user)

	privateKeyFile, _ := ioutil.ReadFile(api.Config.GetString("rsa_private"))
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)

	token := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodRS256.Alg()))
	tokenString, _ := token.SignedString(privateKey)

	return &user, tokenString
}
