package main

import (
	"context"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	alarmrepo "github.com/g0ulartleo/mirante-alerts/internal/alarm/repo"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel/builtins"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	signalrepo "github.com/g0ulartleo/mirante-alerts/internal/signal/repo"
	"github.com/g0ulartleo/mirante-alerts/internal/worker"
	"github.com/hibiken/asynq"
)

func main() {
	alarmRepo, err := alarmrepo.New()
	if err != nil {
		log.Fatalf("Error initializing alarm store: %v", err)
	}
	defer alarmRepo.Close()
	alarmService := alarm.NewAlarmService(alarmRepo)
	err = alarm.InitAlarms(alarmRepo)
	if err != nil {
		log.Fatalf("Error initializing alarm configs: %v", err)
	}

	sentinelFactory := sentinel.NewFactory()
	builtins.Register(sentinelFactory)

	signalRepo, err := signalrepo.New(config.LoadAppConfigFromEnv())
	if err != nil {
		log.Fatalf("Error initializing signal store: %v", err)
	}
	defer signalRepo.Close()
	signalService := signal.NewService(signalRepo)

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
	worker.RegisterTasks(mux, sentinelFactory, signalService, alarmService, asyncClient)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
