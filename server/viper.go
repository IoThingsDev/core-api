package server

func (a API) SetupViper(env string) {
	if env == "prod" {
		a.Config.SetConfigName("conf")
	} else if env == "test" {
		a.Config.SetConfigName("conf.test")
	}
	a.Config.SetConfigType("json")
	a.Config.AddConfigPath("..")
	err := a.Config.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
