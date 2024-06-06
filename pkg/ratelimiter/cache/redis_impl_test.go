package cache

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

var redisInstance *RedisInstance

func TestMain(m *testing.M) {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "localhost:6379", redis.DialConnectTimeout(5*time.Second))
			if err != nil {
				log.Fatalf("Could not connect to Redis: %v", err)
			}
			return conn, err
		},
	}

	redisInstance = NewRedisInstance(pool)

	code := m.Run()
	log.Println("closing connection pool")
	pool.Close()
	os.Exit(code)
}

func TestRedisImplementation(t *testing.T) {

	t.Run("Given a key, when try to increase a value should update properly", func(t *testing.T) {
		expectedKey := "aKey"
		expectedValue := 1

		err := redisInstance.Incr(context.TODO(), expectedKey)
		assert.Nil(t, err)

		value, err := redisInstance.Get(context.TODO(), expectedKey)
		assert.Nil(t, err)
		assert.Equal(t, expectedValue, value)

	})

	t.Run("Given a key, when try to increase value as many time as possible the results should be return properly", func(t *testing.T) {
		expectedKey := "anotherKey"
		expectedValue := rand.Intn(100)

		for range expectedValue {
			err := redisInstance.Incr(context.TODO(), expectedKey)
			assert.Nil(t, err)
		}

		value, err := redisInstance.Get(context.TODO(), expectedKey)
		assert.Nil(t, err)
		assert.Equal(t, expectedValue, value)
	})

	t.Run("Given a key with expiration time when try to get its key after expire should return nil", func(t *testing.T) {
		expectedKey := "oneMoreKey"
		expectedValue := 0
		expectedExpirationTime := 5

		err := redisInstance.Incr(context.TODO(), expectedKey)
		assert.Nil(t, err)
		err = redisInstance.Expire(context.TODO(), expectedKey, expectedExpirationTime)
		assert.Nil(t, err)

		time.Sleep(time.Second * time.Duration(expectedExpirationTime))
		value, err := redisInstance.Get(context.TODO(), expectedKey)

		assert.Equal(t, expectedValue, value)
		assert.NotNil(t, err)
	})
}
