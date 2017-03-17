package server

func (a API) SetupViper(env string) {
	if env == "prod" {
		a.Config.SetConfigName("conf")
		a.Config.SetConfigType("json")
		a.Config.AddConfigPath(".")
		err := a.Config.ReadInConfig()
		if err != nil {
			panic(err)
		}
	} else if env == "test" {
		a.Config.SetConfigName("conf.test")
		a.Config.SetConfigType("json")
		a.Config.AddConfigPath("..")
		err := a.Config.ReadInConfig()
		if err != nil {
			panic(err)
		}
	}
}
