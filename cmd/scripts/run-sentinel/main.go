package main

import (
	"log"
	"os"
	"strconv"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/event_dispatcher/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	conn := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Env().RedisAddr})
	defer conn.Close()

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <sentinel-id>", os.Args[0])
	}

	sentinelID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid sentinel ID: %v", err)
	}

	task, err := tasks.NewSentinelRunTask(sentinelID)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	if _, err := conn.Enqueue(task); err != nil {
		log.Fatalf("Failed to enqueue task: %v", err)
	}

	log.Printf("Task enqueued")
}
