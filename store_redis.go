package oauth2s

import (
	"gopkg.in/oauth2.v4"
	"github.com/go-redis/redis"
	oredis "gopkg.in/go-oauth2/redis.v4"
)

func NewRedisStore(redisOptions ...RedisOptionsSetter) oauth2.TokenStore {
	redisConfig := &redis.Options{}
	for _, fn := range redisOptions {
		fn(redisConfig)
	}
	return oredis.NewRedisStore(redisConfig)
}
