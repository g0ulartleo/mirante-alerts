package tasks

import (
	"context"

	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	backoffice "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/backoffice"
	sentinel "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/sentinel"
	sigTask "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/signal"
	"github.com/hibiken/asynq"
)

func Register(mux *asynq.ServeMux, signalService *signal.Service) {
	mux.HandleFunc(sentinel.TypeSentinelCheckAlert, func(ctx context.Context, task *asynq.Task) error {
		return sentinel.HandleSentinelCheckAlertTask(ctx, task, signalService)
	})
	mux.HandleFunc(sigTask.TypeSignalWrite, func(ctx context.Context, task *asynq.Task) error {
		return sigTask.HandleSignalWriteTask(ctx, task, signalService)
	})
	mux.HandleFunc(backoffice.TypeCleanSignals, func(ctx context.Context, task *asynq.Task) error {
		return backoffice.HandleCleanSignalsTask(ctx, task, signalService)
	})
}
