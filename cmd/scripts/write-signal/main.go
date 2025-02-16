package main

import (
	"log"
	"time"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/event_dispatcher/tasks"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/hibiken/asynq"
)

func main() {
	conn := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Env().RedisAddr})
	defer conn.Close()

	task, err := tasks.NewSignalWriteTask(signal.Signal{
		AlertID:   "project_x-api-google-health-check",
		Status:    signal.StatusHealthy,
		Timestamp: time.Now(),
		Message:   "OK",
	})
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	if _, err := conn.Enqueue(task); err != nil {
		log.Fatalf("Failed to enqueue task: %v", err)
	}

	log.Printf("Task enqueued")
}
