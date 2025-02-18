package main

import (
	"log"
	"os"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/event_dispatcher/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	conn := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Env().RedisAddr})
	defer conn.Close()

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <alert-id>", os.Args[0])
	}

	alertID := os.Args[1]

	task, err := tasks.NewSentinelCheckAlertTask(alertID)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	if _, err := conn.Enqueue(task); err != nil {
		log.Fatalf("Failed to enqueue task: %v", err)
	}

	log.Printf("Task enqueued")
}
