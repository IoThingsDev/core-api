package server

import (
	"os"
	"github.com/joho/godotenv"
)

func (a API) LoadEnvVariables() error {
	var filename string
	env := os.Getenv("BASEAPIENV")

	// Default file is the local one
	if env == "" {
		filename = ".env.local"
	} else if env == "prod" {
		filename = ".env"
	} else {
		filename = ".env."+env
	}

	err := godotenv.Overload(filename)
	return err
}
