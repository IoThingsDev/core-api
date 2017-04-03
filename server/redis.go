package server

import (
	"github.com/dernise/base-api/services"
	"github.com/garyburd/redigo/redis"
)

func (a *API) SetupRedis() {
	pool := &redis.Pool{
		MaxIdle:     a.Config.GetInt("redis_max_idle"),
		MaxActive:   a.Config.GetInt("redis_max_active"), // max number of connections
		IdleTimeout: a.Config.GetDuration("redis_max_timeout"),
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", a.Config.GetString("redis_host"))
		},
	}

	a.Redis = &services.Redis{
		pool,
		a.Config,
	}
}
