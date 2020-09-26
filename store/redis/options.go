package redis

import (
	"crypto/tls"
	"time"

	"github.com/go-redis/redis/v8"
)

type Options struct {
	*redis.Options
	KeyNamespace string
}

type RedisOptionsSetter func(*Options)

func RedisAddr(addr string) RedisOptionsSetter {
	return func(opt *Options) {
		opt.Addr = addr
	}
}

func RedisDB(db int) RedisOptionsSetter {
	return func(opt *Options) {
		opt.DB = db
	}
}

func RedisPassword(password string) RedisOptionsSetter {
	return func(opt *Options) {
		opt.Password = password
	}
}

func RedisUsername(user string) RedisOptionsSetter {
	return func(opt *Options) {
		opt.Username = user
	}
}

func RedisNetwork(network string) RedisOptionsSetter {
	return func(opt *Options) {
		opt.Network = network
	}
}

func RedisMaxRetries(maxRetries int) RedisOptionsSetter {
	return func(opt *Options) {
		opt.MaxRetries = maxRetries
	}
}

func RedisPoolSize(poolSize int) RedisOptionsSetter {
	return func(opt *Options) {
		opt.PoolSize = poolSize
	}
}

func RedisMinIdleConns(minIdleConns int) RedisOptionsSetter {
	return func(opt *Options) {
		opt.MinIdleConns = minIdleConns
	}
}

func RedisMinRetryBackoff(duration time.Duration) RedisOptionsSetter {
	return func(opt *Options) {
		opt.MinRetryBackoff = duration
	}
}

func RedisMaxRetryBackoff(duration time.Duration) RedisOptionsSetter {
	return func(opt *Options) {
		opt.MaxRetryBackoff = duration
	}
}

func RedisDialTimeout(duration time.Duration) RedisOptionsSetter {
	return func(opt *Options) {
		opt.DialTimeout = duration
	}
}

func RedisReadTimeout(duration time.Duration) RedisOptionsSetter {
	return func(opt *Options) {
		opt.DialTimeout = duration
	}
}

func RedisWriteTimeout(duration time.Duration) RedisOptionsSetter {
	return func(opt *Options) {
		opt.WriteTimeout = duration
	}
}

func RedisMaxConnAge(duration time.Duration) RedisOptionsSetter {
	return func(opt *Options) {
		opt.MaxConnAge = duration
	}
}

func RedisPoolTimeout(duration time.Duration) RedisOptionsSetter {
	return func(opt *Options) {
		opt.PoolTimeout = duration
	}
}

func RedisIdleTimeout(duration time.Duration) RedisOptionsSetter {
	return func(opt *Options) {
		opt.IdleTimeout = duration
	}
}

func RedisIdleCheckFrequency(duration time.Duration) RedisOptionsSetter {
	return func(opt *Options) {
		opt.IdleCheckFrequency = duration
	}
}

func RedisLimiter(limiter redis.Limiter) RedisOptionsSetter {
	return func(opt *Options) {
		opt.Limiter = limiter
	}
}

func RedisKeyNamespace(keyNamespace string) RedisOptionsSetter {
	return func(opt *Options) {
		opt.KeyNamespace = keyNamespace
	}
}

func RedisTLSConfig(config *tls.Config) RedisOptionsSetter {
	return func(opt *Options) {
		opt.TLSConfig = config
	}
}
