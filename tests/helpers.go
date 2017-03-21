package tests

import (
	"bytes"
	"github.com/dernise/base-api/server"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"net/http"
	"net/http/httptest"
	"os"
)

type FakeEmailSender struct {
	to          []*mail.Email
	contentType string
	subject     string
	body        string
}

func (f *FakeEmailSender) SendEmail(to []*mail.Email, contentType, subject, body string) (*rest.Response, error) {
	f.to, f.contentType, f.subject, f.body = to, contentType, subject, body
	return &rest.Response{StatusCode: http.StatusOK, Body: "Everything's fine Jean-Miche", Headers: nil}, nil
}

func SetupRouterAndDatabase() *server.API {
	api := server.API{Router: gin.Default(), Config: viper.New()}

	api.LoadEnvVariables()
	api.SetupViper("test")

	session, err := mgo.Dial(os.Getenv("DB_HOST"))
	if err != nil {
		panic(err)
	}

	api.Database = session.DB(os.Getenv("DB_NAME_TEST"))
	api.Database.DropDatabase()

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
