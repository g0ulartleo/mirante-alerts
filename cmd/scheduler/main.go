package main

import (
	"fmt"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	alarmStores "github.com/g0ulartleo/mirante-alerts/internal/alarm/stores"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	alarmTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/alarm"
	backofficeTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/backoffice"
	"github.com/hibiken/asynq"
)

func main() {
	alarmStore, err := alarmStores.NewAlarmStore()
	if err != nil {
		log.Fatalf("Error initializing alarm store: %v", err)
	}
	defer alarmStore.Close()
	alarmService := alarm.NewAlarmService(alarmStore)
	err = alarm.InitAlarms(alarmStore)
	if err != nil {
		log.Fatalf("Error initializing sentinel configs: %v", err)
	}
	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{Addr: config.Env().RedisAddr}, nil)

	alarms, err := alarmService.GetAlarms()
	if err != nil {
		log.Fatalf("Error getting alarms: %v", err)
	}
	for _, sentinelConfig := range alarms {
		task, err := alarmTasks.NewCheckAlarmTask(sentinelConfig.ID)
		if err != nil {
			log.Fatalf("Error creating alarm check task: %v", err)
		}
		cronspec := sentinelConfig.Cron
		if cronspec == "" {
			cronspec = fmt.Sprintf("@every %s", sentinelConfig.Interval)
		}

		entryID, err := scheduler.Register(cronspec, task)
		if err != nil {
			log.Fatalf("Error registering alarm check task: %v", err)
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
