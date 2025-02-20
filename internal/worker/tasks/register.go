package tasks

import (
	"context"

	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	alarm "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/alarm"
	backoffice "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/backoffice"
	sigTask "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/signal"
	"github.com/hibiken/asynq"
)

func Register(mux *asynq.ServeMux, signalService *signal.Service, asyncClient *asynq.Client) {
	mux.HandleFunc(alarm.TypeCheckAlarm, func(ctx context.Context, task *asynq.Task) error {
		return alarm.HandleCheckAlarmTask(ctx, task, signalService, asyncClient)
	})
	mux.HandleFunc(sigTask.TypeSignalWrite, func(ctx context.Context, task *asynq.Task) error {
		return sigTask.HandleSignalWriteTask(ctx, task, signalService)
	})
	mux.HandleFunc(backoffice.TypeCleanSignals, func(ctx context.Context, task *asynq.Task) error {
		return backoffice.HandleCleanSignalsTask(ctx, task, signalService)
	})
	mux.HandleFunc(alarm.TypeNotify, func(ctx context.Context, task *asynq.Task) error {
		return alarm.HandleNotifyTask(ctx, task)
	})
}
