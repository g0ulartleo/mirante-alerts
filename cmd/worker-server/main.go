package main

import (
	"context"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	alarmfactory "github.com/g0ulartleo/mirante-alerts/internal/alarm/factory"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel/builtin"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	signalfactory "github.com/g0ulartleo/mirante-alerts/internal/signal/factory"
	"github.com/g0ulartleo/mirante-alerts/internal/worker"
	"github.com/hibiken/asynq"
)

func main() {
	alarmStore, err := alarmfactory.New()
	if err != nil {
		log.Fatalf("Error initializing alarm store: %v", err)
	}
	defer alarmStore.Close()
	alarmService := alarm.NewAlarmService(alarmStore)
	err = alarm.InitAlarms(alarmStore)
	if err != nil {
		log.Fatalf("Error initializing alarm configs: %v", err)
	}

	builtin.Register(sentinel.Factory)

	signalStore, err := signalfactory.New(config.LoadAppConfigFromEnv())
	if err != nil {
		log.Fatalf("Error initializing signal store: %v", err)
	}
	defer signalStore.Close()
	signalService := signal.NewService(signalStore)

	asyncClient := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Env().RedisAddr})
	defer asyncClient.Close()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.Env().RedisAddr},
		asynq.Config{
			Concurrency: 10,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Printf("Error processing task %s: %v", task.Type(), err)
			}),
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	worker.RegisterTasks(mux, signalService, alarmService, asyncClient)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
