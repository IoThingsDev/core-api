package server

import mgo "gopkg.in/mgo.v2"

func (a *API) SetupDatabase() (*mgo.Session, error) {
	session, err := mgo.Dial(a.Config.GetString("db_host"))
	if err != nil {
		return nil, err
	}

	a.Database = session.DB(a.Config.GetString("db_name"))

	return session, nil
}
