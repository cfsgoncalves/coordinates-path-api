package repository

import (
	"context"
	"fmt"
	"meight/configuration"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var ctx = context.Background()

type Redis struct {
	Redis redis.Client
}

func NewRedis() *Redis {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configuration.GetEnvAsString("REDIS_SERVER", "localhost"), configuration.GetEnvAsString("REDIS_PORT", "6379")), // use default Addr
		Password: configuration.GetEnvAsString("REDIS_PASSWORD", ""),                                                                                  // no password set
		DB:       configuration.GetEnvAsInt("REDIS_DB", 0),                                                                                            // use default DB
	})

	status, err := redisClient.Ping(ctx).Result()

	if err != nil {
		log.Error().Msgf("repository.NewRedis(): Error yield trying to acess redis client. Error: %s", err)
		return &Redis{}
	}

	if status != "PONG" {
		log.Error().Msgf("repository.NewRedis(): Error while trying to acess redis client. Status is %s", status)
		return &Redis{}
	}

	return &Redis{Redis: *redisClient}
}

func (r *Redis) Insert(key string, value string) error {
	TTL := time.Duration(configuration.GetEnvAsInt("REDIS_TTL", 15*60)) * time.Second

	err := r.Redis.Set(ctx, key, value, TTL).Err()
	if err != nil {
		log.Error().Msgf("repository.Get(): Error while inserting data into reddis for key %s. Error %s", key, err)
		return err
	}
	return nil
}

func (r *Redis) Get(key string) (string, error) {

	val, err := r.Redis.Get(ctx, key).Result()

	if err != nil && err.Error() == "redis: nil" {
		log.Debug().Msgf("repository.Get(): No value found for key %s", key)
	}

	if err != nil && err.Error() != "redis: nil" {
		log.Error().Msgf("repository.Get(): Error while fetching data from reddis for %s. Error %s", key, err)
		return "", err
	}

	return val, nil

}

func (r *Redis) Ping() bool {
	status, err := r.Redis.Ping(ctx).Result()

	if err != nil {
		log.Error().Msgf("repository.Ping(): Error yield trying to acess redis client. Error: %s", err)
		return false
	}

	if status != "PONG" {
		log.Error().Msgf("repository.Ping(): Error while trying to acess redis client. Status is %s", status)
		return false
	}

	return true
}
