package redis

import (
	"github.com/go-redis/redis"
	oredis "gopkg.in/go-oauth2/redis.v4"
	"gopkg.in/oauth2.v4"
)

func New(redisOptions ...RedisOptionsSetter) oauth2.TokenStore {
	redisConfig := &Options{
		Options: &redis.Options{},
	}
	for _, fn := range redisOptions {
		fn(redisConfig)
	}
	var keyNamespaces []string
	if len(redisConfig.KeyNamespace) > 0 {
		keyNamespaces = append(keyNamespaces, redisConfig.KeyNamespace)
	}
	return oredis.NewRedisStore(redisConfig.Options, keyNamespaces...)
}

func NewWithClient(client *redis.Client, keyNamespaces ...string) oauth2.TokenStore {
	return oredis.NewRedisStoreWithCli(client, keyNamespaces...)
}
