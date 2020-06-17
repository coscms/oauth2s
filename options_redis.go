package oauth2s

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisOptionsSetter func(*redis.Option)

func RedisAddr(addr string) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.Addr = addr
	}
}

func RedisDB(db int) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.DB = db
	}
}

func RedisPassword(password string) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.Password = password
	}
}

func RedisUsername(user string) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.Username = user
	}
}

func RedisNetwork(network string) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.Network = network
	}
}

func RedisMaxRetries(maxRetries int) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.MaxRetries = maxRetries
	}
}

func RedisPoolSize(poolSize int) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.PoolSize = poolSize
	}
}

func RedisMinIdleConns(minIdleConns int) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.MinIdleConns = minIdleConns
	}
}

func RedisMinRetryBackoff(duration time.Duration) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.MinRetryBackoff = duration
	}
}

func RedisMaxRetryBackoff(duration time.Duration) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.MaxRetryBackoff = duration
	}
}

func RedisDialTimeout(duration time.Duration) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.DialTimeout = duration
	}
}

func RedisReadTimeout(duration time.Duration) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.DialTimeout = duration
	}
}

func RedisWriteTimeout(duration time.Duration) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.WriteTimeout = duration
	}
}

func RedisMaxConnAge(duration time.Duration) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.MaxConnAge = duration
	}
}

func RedisPoolTimeout(duration time.Duration) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.PoolTimeout = duration
	}
}

func RedisIdleTimeout(duration time.Duration) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.IdleTimeout = duration
	}
}

func RedisIdleCheckFrequency(duration time.Duration) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.IdleCheckFrequency = duration
	}
}

func RedisLimiter(limiter redis.Limiter) RedisOptionsSetter {
	return func(opt *redis.Option) {
		opt.Limiter = limiter
	}
}
