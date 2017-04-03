package middlewares

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dernise/base-api/helpers"
	"github.com/dernise/base-api/services"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
)

func RateMiddleware(redis *services.Redis, config *viper.Viper) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn := redis.Pool.Get()
		defer conn.Close()

		if !config.GetBool("rate_limit_activated") {
			return
		}

		unixTime := int64(time.Now().Unix())
		keyName := c.ClientIP() + ":" + strconv.FormatInt(unixTime, 10)

		var count int = -1
		data, err := conn.Do("GET", keyName)
		if err != nil {
			return
		}

		if data != nil {
			if err := json.Unmarshal(data.([]byte), &count); err != nil {
				return
			}
		}

		if count != -1 && count >= config.GetInt("rate_limit_requests_per_second") {
			c.AbortWithError(http.StatusTooManyRequests, helpers.ErrorWithCode("too_many_requests", "You sent too many requests over the last second."))
		} else {
			conn.Do("INCR", keyName)
			conn.Do("EXPIRE", keyName, 10)
		}
	}
}
