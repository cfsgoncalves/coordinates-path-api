package repository

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisInsert(t *testing.T) {

	t.Run("happy_path_integration", func(t *testing.T) {

		redis := NewRedis()

		redis.Insert("2024-06-24", "fooo")

		bar, _ := redis.Get("2024-06-24")

		assert.Equal(t, bar, "fooo")
	})

	t.Run("test_ttl_integration", func(t *testing.T) {
		os.Setenv("REDIS_TTL", "1")

		redis := NewRedis()

		redis.Insert("2024-06-24", "fooo")

		time.Sleep(2 * time.Second)

		bar, _ := redis.Get("2024-06-24")

		assert.NotEqual(t, bar, "fooo")

	})

	t.Run("fail_to_insert", func(t *testing.T) {
		redis := NewRedis()

		db, mock := redismock.NewClientMock()

		redis.Redis = *db

		mock.ExpectSet("2024-06-24", "fooo", 1).SetErr(errors.New("error"))

		err := redis.Insert("2024-06-24", "fooo")

		assert.Error(t, err, errors.New("error"))
		assert.NotNil(t, err, "This should be not nill")
	})

}

func TestRedisGet(t *testing.T) {

	t.Run("fail_to_get", func(t *testing.T) {
		redis := NewRedis()

		db, mock := redismock.NewClientMock()

		redis.Redis = *db

		mock.ExpectGet("2024-06-24").SetErr(errors.New("error"))

		_, err := redis.Get("2024-06-24")

		assert.Error(t, err, errors.New("error"))
		assert.NotNil(t, err, "This should be not nill")
	})
}
