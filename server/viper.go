package server

import (
	"os"

	"github.com/joho/godotenv"
)

func (a *API) SetupViper() error {

	filename := ".env"
	switch os.Getenv("BASEAPI_ENV") {
	case "testing":
		filename = "../.env.testing"
	case "prod":
		filename = ".env.prod"
	}

	err := godotenv.Overload(filename)
	if err != nil {
		return err
	}

	a.Config.SetEnvPrefix("baseapi")
	a.Config.AutomaticEnv()

	return nil
}
