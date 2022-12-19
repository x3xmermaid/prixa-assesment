package redis_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	nredis "prixa-assesment/app/db/redis"

	redis "github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

func TestNewRedis_NOK(t *testing.T) {
	_ = redigomock.NewConn()

	t.Run("TestNewRedis_OK", func(t *testing.T) {
		network := "tcp"
		address := "localhost:6379"
		password := "pass"
		timeout := int64(1)

		redisdb, err := nredis.NewRedis(network, address, password, timeout, 300)
		if err != nil {
			t.Errorf("Connection should not be error")
		}

		redisdb.GetValue(rediskey)
	})
}

const rediskey = "GetAllLinkUtilization"

func TestIsAvailable(t *testing.T) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		MaxIdle:     1,
		IdleTimeout: 1 * time.Second,
		MaxActive:   1,
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	t.Run("TestIsAvailable_OK", func(t *testing.T) {
		status := "true"
		byteStatus, _ := json.Marshal(status)
		conn.Command("GET", rediskey).Expect(byteStatus)

		redisdb := nredis.Redis{
			Pool:    pool,
			Timeout: 1 * time.Second,
		}

		if ok := redisdb.IsAvailable(rediskey); !ok {
			t.Errorf("TestIsAvailable_OK should be %v but found: %v", true, ok)
		}
	})

	t.Run("TestIsAvailable_NOK", func(t *testing.T) {
		conn.Command("GET", rediskey).ExpectError(fmt.Errorf("Intentionally error"))

		redisdb := nredis.Redis{
			Pool:    pool,
			Timeout: 1 * time.Second,
		}

		if ok := redisdb.IsAvailable(rediskey); ok {
			t.Errorf("TestIsAvailable_NOK should be error")
		}
	})

	t.Run("TestIsAvailable_NOK_Nil", func(t *testing.T) {
		conn.Command("GET", rediskey).ExpectError(redis.ErrNil)

		redisdb := nredis.Redis{
			Pool:    pool,
			Timeout: 1 * time.Second,
		}

		if ok := redisdb.IsAvailable(rediskey); ok {
			t.Errorf("TestIsAvailable_NOK should be error")
		}
	})

	t.Run("TestIsAvailable_NOK_PoolGet", func(t *testing.T) {
		conn.Command("GET", rediskey).ExpectError(redis.ErrNil)
		pool.Close()
		redisdb := nredis.Redis{
			Pool:    pool,
			Timeout: 1 * time.Second,
		}

		if ok := redisdb.IsAvailable(rediskey); ok {
			t.Errorf("TestIsAvailable_NOK_GetPool should be error: %v", ok)
		}

	})
}

func TestPut(t *testing.T) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		MaxIdle:   2,
		MaxActive: 20,
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
	}

	status := 21
	byteStatus, _ := json.Marshal(status)

	t.Run("TestPut_NOK_Marshal", func(t *testing.T) {
		redisdb := nredis.Redis{
			Pool:          pool,
			Timeout:       1 * time.Second,
			CacheDuration: 1000,
		}

		invalidJSON := make(chan int)

		conn.Command("SET", rediskey, byteStatus)

		err := redisdb.Put(rediskey, invalidJSON)
		if err == nil {
			t.Errorf("TestPut_OK should be error but found")
		}
	})

	t.Run("TestPut_OK", func(t *testing.T) {

		redisdb := nredis.Redis{
			Pool:          pool,
			Timeout:       1 * time.Second,
			CacheDuration: 1000,
		}

		conn.Command("SET", rediskey, byteStatus)

		err := redisdb.Put(rediskey, status)
		if err != nil {
			t.Errorf("TestPut_OK should not be error but found: %v", err.Error())
		}
	})

	t.Run("TestPut_NOK", func(t *testing.T) {
		conn.Command("SET", rediskey, byteStatus).ExpectError(fmt.Errorf("Intentionally error"))

		redisdb := nredis.Redis{
			Pool:    pool,
			Timeout: 1 * time.Second,
		}

		err := redisdb.Put(rediskey, status)
		if err == nil {
			t.Errorf("TestPut_OK should be error but found")
		}
	})

	t.Run("TestPut_NOK_Nil", func(t *testing.T) {
		conn.Command("SET", rediskey, byteStatus).ExpectError(redis.ErrNil)

		redisdb := nredis.Redis{
			Pool:    pool,
			Timeout: 1 * time.Second,
		}

		err := redisdb.Put(rediskey, status)
		if err == nil {
			t.Errorf("TestPut_OK should be error but found")
		}
	})

	t.Run("TestPut_NOK_PoolGet", func(t *testing.T) {
		conn.Command("SET", rediskey, byteStatus).ExpectError(redis.ErrNil)
		pool.Close()
		redisdb := nredis.Redis{
			Pool:    pool,
			Timeout: 1 * time.Second,
		}

		err := redisdb.Put(rediskey, status)
		if err == nil {
			t.Errorf("TestPut_OK should be error but found")
		}
	})
}

func TestGetValue(t *testing.T) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		MaxIdle:   2,
		MaxActive: 20,
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
	}

	t.Run("TestGetValue_OK", func(t *testing.T) {
		status := 21
		byteStatus, _ := json.Marshal(status)
		conn.Command("GET", rediskey).Expect(byteStatus)

		redisdb := nredis.Redis{
			Pool:    pool,
			Timeout: 1 * time.Second,
		}

		_, err := redisdb.GetValue(rediskey)
		if err != nil {
			t.Errorf("TestGetValue_OK should not be error but found: %v", err.Error())
		}
	})

	t.Run("TestGetValue_NOK", func(t *testing.T) {
		conn.Command("GET", rediskey).ExpectError(fmt.Errorf("Intentionally error"))

		redisdb := nredis.Redis{
			Pool:    pool,
			Timeout: 1 * time.Second,
		}

		lastUpdated, err := redisdb.GetValue(rediskey)
		if err == nil {
			t.Errorf("TestGetValue_OK should be error but found: %v", lastUpdated)
		}
	})

	t.Run("TestGetValue_NOK_Nil", func(t *testing.T) {
		conn.Command("GET", rediskey).ExpectError(redis.ErrNil)

		redisdb := nredis.Redis{
			Pool:    pool,
			Timeout: 1 * time.Second,
		}

		lastUpdated, err := redisdb.GetValue(rediskey)
		if err == nil {
			t.Errorf("TestGetValue_OK should be error but found: %v", lastUpdated)
		}
	})

	t.Run("TestGetValue_NOK_PoolGet", func(t *testing.T) {
		conn.Command("GET", rediskey).ExpectError(redis.ErrNil)
		pool.Close()
		redisdb := nredis.Redis{
			Pool:    pool,
			Timeout: 1 * time.Second,
		}

		lastUpdated, err := redisdb.GetValue(rediskey)
		if err == nil {
			t.Errorf("TestGetValue_OK should be error but found: %v", lastUpdated)
		}
	})
}
