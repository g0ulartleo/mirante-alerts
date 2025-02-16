package main

import (
	"fmt"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/event_dispatcher/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	err := config.InitAlerts()
	if err != nil {
		log.Fatalf("Error initializing sentinel configs: %v", err)
	}
	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{Addr: config.Env().RedisAddr}, nil)

	for _, sentinelConfig := range config.Alerts {
		task, err := tasks.NewSentinelCheckAlertTask(sentinelConfig.ID)
		if err != nil {
			log.Fatalf("Error creating sentinel check alert task: %v", err)
		}
		cronspec := sentinelConfig.Cron
		if cronspec == "" {
			cronspec = fmt.Sprintf("@every %s", sentinelConfig.Interval)
		}

		entryID, err := scheduler.Register(cronspec, task)
		if err != nil {
			log.Fatalf("Error registering sentinel check alert task: %v", err)
		}
		log.Printf("registered an entry: %q\n", entryID)
	}

	cleanSignalsTask, err := tasks.NewCleanSignalsTask()
	if err != nil {
		log.Fatalf("Error creating clean signals task: %v", err)
	}
	scheduler.Register("@every 1d", cleanSignalsTask)

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}

}
