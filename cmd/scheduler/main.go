package main

import (
	"fmt"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/event_dispatcher/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	err := config.InitSentinelConfigs()
	if err != nil {
		log.Fatalf("Error initializing sentinel configs: %v", err)
	}
	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{Addr: config.Env().RedisAddr}, nil)

	for _, sentinelConfig := range config.SentinelConfigs {
		task, err := tasks.NewSentinelRunTask(sentinelConfig.ID)
		if err != nil {
			log.Fatalf("Error creating sentinel run task: %v", err)
		}
		cronspec := sentinelConfig.Cron
		if cronspec == "" {
			cronspec = fmt.Sprintf("@every %s", sentinelConfig.Interval)
		}

		entryID, err := scheduler.Register(cronspec, task)
		if err != nil {
			log.Fatalf("Error registering sentinel run task: %v", err)
		}
		log.Printf("registered an entry: %q\n", entryID)
	}

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}

}
