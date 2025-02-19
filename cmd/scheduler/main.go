package main

import (
	"fmt"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	backofficeTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/backoffice"
	sentinelTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/sentinel"
	"github.com/hibiken/asynq"
)

func main() {
	err := alarm.InitAlarms()
	if err != nil {
		log.Fatalf("Error initializing sentinel configs: %v", err)
	}
	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{Addr: config.Env().RedisAddr}, nil)

	for _, sentinelConfig := range alarm.Alarms {
		task, err := sentinelTasks.NewSentinelCheckAlarmTask(sentinelConfig.ID)
		if err != nil {
			log.Fatalf("Error creating sentinel check alarm task: %v", err)
		}
		cronspec := sentinelConfig.Cron
		if cronspec == "" {
			cronspec = fmt.Sprintf("@every %s", sentinelConfig.Interval)
		}

		entryID, err := scheduler.Register(cronspec, task)
		if err != nil {
			log.Fatalf("Error registering sentinel check alarm task: %v", err)
		}
		log.Printf("registered an entry: %q\n", entryID)
	}

	cleanSignalsTask, err := backofficeTasks.NewCleanSignalsTask()
	if err != nil {
		log.Fatalf("Error creating clean signals task: %v", err)
	}
	scheduler.Register("@every 1d", cleanSignalsTask)

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}

}
