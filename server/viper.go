package server

func (a API) SetupViper() {
	a.Config.SetConfigName("conf")
	a.Config.SetConfigType("json")
	a.Config.AddConfigPath(".")
	err := a.Config.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
