package server

func (a API) SetupViper(env string) error {
	if env == "prod" {
		a.Config.SetConfigName("conf")
		a.Config.AddConfigPath(".")
	} else if env == "test" {
		a.Config.SetConfigName("conf.test")
		a.Config.AddConfigPath("..") // TODO: REFACTOR THIS
	}
	a.Config.SetConfigType("json")
	err := a.Config.ReadInConfig()
	return err
}
