package server

import (
	"os"

	"time"

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

	a.SetupViperDefaults()

	return nil
}

func (a *API) SetupViperDefaults() {
	a.Config.SetDefault("redis_max_idle", 80)
	a.Config.SetDefault("redis_max_active", 12000)
	a.Config.SetDefault("redis_max_timeout", 240*time.Second)
	a.Config.SetDefault("redis_cache_expiration", 10)
	a.Config.SetDefault("rate_limit_requests_per_second", 5)
	a.Config.SetDefault("rate_limit_activated", true)
}
