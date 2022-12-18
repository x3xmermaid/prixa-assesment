package redis

import (
	"time"

	redis "github.com/gomodule/redigo/redis"
)

const (
	maxIdle   = 80
	maxActive = 12000
)

// Redis implements IRedis
type Redis struct {
	Pool          *redis.Pool
	Timeout       time.Duration
	CacheDuration int64
}

// NewRedis construct Redis
func NewRedis(network, address, password string, timeout int64, cacheDuration int64) (*Redis, error) {
	redis := &Redis{}

	return redis, nil
}
