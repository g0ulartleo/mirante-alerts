package main

import (
	"fmt"
	"log"
	"time"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	alarmStores "github.com/g0ulartleo/mirante-alerts/internal/alarm/stores"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	alarmTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/alarm"
	backofficeTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/backoffice"
	"github.com/hibiken/asynq"
)

type AlarmConfigProvider struct {
	alarmService *alarm.AlarmService
}

func (p *AlarmConfigProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	alarms, err := p.alarmService.GetAlarms()
	if err != nil {
		return nil, fmt.Errorf("error getting alarms: %v", err)
	}

	var configs []*asynq.PeriodicTaskConfig
	for _, alarmConfig := range alarms {
		task, err := alarmTasks.NewCheckAlarmTask(alarmConfig.ID)
		if err != nil {
			return nil, fmt.Errorf("error creating alarm check task: %v", err)
		}
		cronspec := alarmConfig.Cron
		if cronspec == "" {
			cronspec = fmt.Sprintf("@every %s", alarmConfig.Interval)
		}
		configs = append(configs, &asynq.PeriodicTaskConfig{
			Cronspec: cronspec,
			Task:     task,
		})
	}

	cleanSignalsTask, err := backofficeTasks.NewCleanSignalsTask()
	if err != nil {
		return nil, fmt.Errorf("error creating clean signals task: %v", err)
	}
	configs = append(configs, &asynq.PeriodicTaskConfig{
		Cronspec: "@every 24h",
		Task:     cleanSignalsTask,
	})

	return configs, nil
}

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

	provider := &AlarmConfigProvider{
		alarmService: alarmService,
	}

	mgr, err := asynq.NewPeriodicTaskManager(
		asynq.PeriodicTaskManagerOpts{
			RedisConnOpt:               asynq.RedisClientOpt{Addr: config.Env().RedisAddr},
			PeriodicTaskConfigProvider: provider,
			SyncInterval:               30 * time.Second,
		},
	)
	if err != nil {
		log.Fatalf("Error creating periodic task manager: %v", err)
	}

	if err := mgr.Run(); err != nil {
		log.Fatalf("Error running periodic task manager: %v", err)
	}
}
