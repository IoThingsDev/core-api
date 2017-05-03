package config

import "github.com/spf13/viper"

type conf struct {
	*viper.Viper
}

func New(viper *viper.Viper) *conf {
	return &conf{viper}
}
