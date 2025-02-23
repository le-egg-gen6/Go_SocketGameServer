package main

import (
	"go_game_server/cache/redis"
	"go_game_server/config"
	"go_game_server/utils/log"
)

func main() {
	log.InitializeLogger()
	defer CleanupUnfinishedTask()

	cfg := config.LoadConfig()

	//Redis
	redis.InitializeRedisInstance(cfg)

}

func CleanupUnfinishedTask() {
	redis.CloseConnection()
	log.CleanupQueuedLogs()
}
