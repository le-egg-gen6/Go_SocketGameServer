package main

import (
	"go_game_server/config"
	"go_game_server/utils/log"
)

func main() {
	log.InitializeLogger()
	defer CleanupUnfinishedTask()

	cfg := config.LoadConfig()

	log.LogInstance.Info(cfg.DB_PASSWORD)

}

func CleanupUnfinishedTask() {
	log.CleanupQueuedLogs()
}
