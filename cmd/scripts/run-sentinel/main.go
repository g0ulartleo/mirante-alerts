package main

import (
	"log"
	"os"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	sentinelTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/sentinel"
	"github.com/hibiken/asynq"
)

func main() {
	conn := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Env().RedisAddr})
	defer conn.Close()

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <alarm-id>", os.Args[0])
	}

	alarmID := os.Args[1]

	task, err := sentinelTasks.NewSentinelCheckAlarmTask(alarmID)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	if _, err := conn.Enqueue(task); err != nil {
		log.Fatalf("Failed to enqueue task: %v", err)
	}

	log.Printf("Task enqueued")
}
