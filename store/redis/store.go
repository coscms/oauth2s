package redis

import (
	"github.com/go-redis/redis"
	oredis "gopkg.in/go-oauth2/redis.v4"
	"gopkg.in/oauth2.v4"
)

func New(redisOptions ...RedisOptionsSetter) oauth2.TokenStore {
	redisConfig := &redis.Options{}
	for _, fn := range redisOptions {
		fn(redisConfig)
	}
	return oredis.NewRedisStore(redisConfig)
}
