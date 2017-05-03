package services

import (
	"encoding/json"
	"errors"

	"context"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
)

type Redis struct {
	Pool   *redis.Pool
	Config *viper.Viper
}

func GetRedis(c context.Context) *Redis {
	return c.Value("redis").(*Redis)
}

// Gets the object from the redis store, if not found, returns an err
func (r *Redis) GetValueForKey(key string, result interface{}) error {
	redisConn := r.Pool.Get()
	defer redisConn.Close()

	data, err := redisConn.Do("GET", key)
	if err != nil {
		return err
	}

	if data == nil {
		return errors.New("The given key was not found in the redis store")
	}

	// Unmarshalls the json
	if err := json.Unmarshal(data.([]byte), &result); err != nil {
		return err
	}

	return nil
}

func (r *Redis) SetValueForKey(key string, value interface{}) error {
	redisConn := r.Pool.Get()
	defer redisConn.Close()

	// Marshals the value to JSON format
	marshaledData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Sets the value to the redis keystore
	if _, err := redisConn.Do("SET", key, marshaledData); err != nil {
		return err
	}

	// Sets the expiration
	if _, err := redisConn.Do("EXPIRE", key, r.Config.GetInt("redis_cache_expiration")*60); err != nil {
		return err
	}

	return nil
}

// Invalidates an object in the redis store.
func (r *Redis) InvalidateObject(id string) error {
	redisConn := r.Pool.Get()
	defer redisConn.Close()

	if _, err := redisConn.Do("DEL", id); err != nil {
		return err
	}

	return nil
}

func (r *Redis) UpdateEmailRateLimit(ip string) error {
	redisConn := r.Pool.Get()
	defer redisConn.Close()

	fullKeyName := ip + ":email"
	count := 1

	err := r.GetValueForKey(fullKeyName, &count)
	if err != nil {
		r.SetValueForKey(fullKeyName, &count)
		return nil
	}

	if _, err := redisConn.Do("INCR", fullKeyName); err != nil {
		return err
	}

	if count > 10 {
		return errors.New("You reached the maximum rate on this request.")
	}
	return nil
}
