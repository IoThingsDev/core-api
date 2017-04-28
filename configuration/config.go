package configuration

import "github.com/spf13/viper"

type config struct {
	*viper.Viper
}

func New(viper *viper.Viper) *config {
	return &config{viper}
}
