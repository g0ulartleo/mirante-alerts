package main

import (
	"context"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel/builtin"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel/custom"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/g0ulartleo/mirante-alerts/internal/signal/stores"
	"github.com/g0ulartleo/mirante-alerts/internal/worker/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	err := alarm.InitAlarms()
	if err != nil {
		log.Fatalf("Error initializing alarm configs: %v", err)
	}

	builtin.Register(sentinel.Factory)
	custom.Register(sentinel.Factory)

	signalStore, err := stores.NewStore(config.LoadAppConfigFromEnv())
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
	tasks.Register(mux, signalService)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
