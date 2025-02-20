package main

import (
	"log"
	"os"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	alarmTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/alarm"
	"github.com/hibiken/asynq"
)

func main() {
	conn := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Env().RedisAddr})
	defer conn.Close()

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <alarm-id>", os.Args[0])
	}

	alarmID := os.Args[1]

	task, err := alarmTasks.NewCheckAlarmTask(alarmID)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	if _, err := conn.Enqueue(task); err != nil {
		log.Fatalf("Failed to enqueue task: %v", err)
	}

	log.Printf("Task enqueued")
}
