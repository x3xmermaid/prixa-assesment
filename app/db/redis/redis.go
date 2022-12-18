package redis

import (
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
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
	timeoutDuration := time.Duration(timeout) * time.Second
	timeoutOption := redis.DialConnectTimeout(timeoutDuration)
	passwordOption := redis.DialPassword(password)
	pool := &redis.Pool{
		MaxIdle:   maxIdle,
		MaxActive: maxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(network, address, timeoutOption, passwordOption)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	redis := &Redis{
		Pool:          pool,
		Timeout:       timeoutDuration,
		CacheDuration: cacheDuration,
	}

	return redis, nil
}

// IsAvailable check wether a cache is available during the time it was called
func (r *Redis) IsAvailable(key string) bool {
	conn := r.Pool.Get()
	if conn.Err() != nil {
		return false
	}
	defer conn.Close()

	_, err := redis.String(redis.DoWithTimeout(conn, r.Timeout, "GET", key))
	if err == redis.ErrNil {
		logrus.Println("Key not found")
		return false
	} else if err != nil {
		return false
	}

	return true
}

// Put set new value to the redis
func (r *Redis) Put(key string, value interface{}) error {
	result, err := json.Marshal(value)
	if err != nil {
		return err
	}

	conn := r.Pool.Get()
	if conn.Err() != nil {
		return conn.Err()
	}
	defer conn.Close()

	_, err = redis.DoWithTimeout(conn, r.Timeout, "SET", key, result)
	if err != nil {
		return err
	}

	return nil
}

// GetValue retrieves the cache value from redis
func (r *Redis) GetValue(key string) ([]byte, error) {
	conn := r.Pool.Get()
	if conn.Err() != nil {
		return nil, conn.Err()
	}
	defer conn.Close()

	data, err := redis.Bytes(redis.DoWithTimeout(conn, r.Timeout, "GET", key))
	if err == redis.ErrNil {
		logrus.Println("Key not found")
		return nil, fmt.Errorf("key not found")
	} else if err != nil {
		return nil, err
	}

	return data, nil
}
