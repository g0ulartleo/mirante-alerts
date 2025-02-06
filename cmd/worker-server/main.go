package main

import (
	"context"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/config"
	"github.com/g0ulartleo/mirante-alerts/internal/event_dispatcher/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	config.InitSentinelConfigs()
	config.InitSentinelFactory()

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
	mux.HandleFunc(tasks.TypeSignalWrite, tasks.HandleSignalWriteTask)
	mux.HandleFunc(tasks.TypeSentinelRun, tasks.HandleSentinelRunTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
