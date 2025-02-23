package redis

import "go_game_server/config"

type RedisConfig struct {
	REDIS_HOST     string
	REDIS_PORT     int
	REDIS_PASSWORD string
	REDIS_DATABASE int
}

func LoadRedisConfig(cfg *config.Config) *RedisConfig {
	return &RedisConfig{
		REDIS_HOST:     cfg.REDIS_HOST,
		REDIS_PORT:     cfg.REDIS_PORT,
		REDIS_PASSWORD: cfg.REDIS_PASSWORD,
		REDIS_DATABASE: cfg.REDIS_DATABASE,
	}
}
