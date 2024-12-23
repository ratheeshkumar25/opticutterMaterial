package config

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisService represents the client of redis.
type RedisService struct {
	Client *redis.Client
}

// SetupRedis initialisez redis server with configuration variables.
func SetupRedis(cfg *Config) (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.REDISHOST,
		DB:   0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.New("failed to connect to redis")
	}

	return &RedisService{
		Client: client,
	}, nil
}

// SetDataInRedis will set  data in redis with a key and expiration time.
func (r *RedisService) SetDataInRedis(key string, value []byte, expTime time.Duration) error {
	err := r.Client.Set(context.Background(), key, value, expTime).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetFromRedis will help to retrieve the data from redis.
func (r *RedisService) GetFromRedis(key string) (string, error) {
	jsonData, err := r.Client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return jsonData, nil
}

// DeleteDataFromRedis will delete data from Redis for a given key.
func (r *RedisService) DeleteDataFromRedis(key string) error {
	ctx := context.Background()
	err := r.Client.Del(ctx, key).Err() 
	if err != nil {
		return err // Return the error if deletion fails
	}
	return nil
}
