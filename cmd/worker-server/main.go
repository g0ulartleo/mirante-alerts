package main

import (
	"context"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/event_dispatcher/tasks"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/g0ulartleo/mirante-alerts/internal/signal/stores"
	"github.com/hibiken/asynq"
)

func main() {
	err := config.InitAlerts()
	if err != nil {
		log.Fatalf("Error initializing alert configs: %v", err)
	}
	config.InitSentinelFactory()

	signalStore, err := stores.NewStore(config.LoadSignalsDatabaseConfigFromEnv())
	if err != nil {
		log.Fatalf("Error initializing signal store: %v", err)
	}
	defer signalStore.Close()
	signalService := signal.NewService(signalStore)

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
	mux.HandleFunc(tasks.TypeSignalWrite, func(ctx context.Context, t *asynq.Task) error {
		return tasks.HandleSignalWriteTask(ctx, t, signalService)
	})
	mux.HandleFunc(tasks.TypeSentinelCheckAlert, func(ctx context.Context, t *asynq.Task) error {
		return tasks.HandleSentinelCheckAlertTask(ctx, t, signalService)
	})
	mux.HandleFunc(tasks.TypeCleanSignals, func(ctx context.Context, t *asynq.Task) error {
		return tasks.HandleCleanSignalsTask(ctx, t, signalService)
	})

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
