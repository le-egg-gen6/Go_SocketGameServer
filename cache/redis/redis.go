package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go_game_server/config"
	"time"
)

type Client struct {
	client *redis.Client
	ctx    context.Context
}

var Instance *Client

func InitializeNewRedisClient(redisConfig *RedisConfig) (*Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.REDIS_HOST, redisConfig.REDIS_PORT),
		Password: redisConfig.REDIS_PASSWORD,
		DB:       redisConfig.REDIS_DATABASE,
	})

	ctx := context.Background()

	_, err := redisClient.Ping(ctx).Result()

	if err != nil {
		return nil, err
	}

	return &Client{
		client: redisClient,
		ctx:    ctx,
	}, nil
}

func InitializeRedisInstance(cfg *config.Config) {
	redisCfg := LoadRedisConfig(cfg)
	instance, err := InitializeNewRedisClient(redisCfg)

	if err != nil {
		panic("Redis client failed to initialize " + err.Error())
	}

	Instance = instance
}

func (r *Client) SetInteger(key string, value int64, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *Client) GetInteger(key string) (int64, error) {
	result, err := r.client.Get(r.ctx, key).Int64()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return result, err
}

func (r *Client) SetString(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *Client) GetString(key string) (string, error) {
	result, err := r.client.Get(r.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return result, err
}

func (r *Client) SetDouble(key string, value float64, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *Client) GetDouble(key string) (float64, error) {
	result, err := r.client.Get(r.ctx, key).Float64()
	if errors.Is(err, redis.Nil) {
		return 0.0, nil
	}
	return result, err
}

func (r *Client) Del(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func CloseConnection() {
	if Instance != nil {
		err := Instance.client.Close()
		if err != nil {
			// log
		}
	}
}
