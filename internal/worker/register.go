package worker

import (
	"context"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/g0ulartleo/mirante-alerts/internal/worker/tasks"
	"github.com/hibiken/asynq"
)

func RegisterTasks(mux *asynq.ServeMux, signalService *signal.Service, alarmService *alarm.AlarmService, asyncClient *asynq.Client) {
	mux.HandleFunc(tasks.TypeAlarmCheck, func(ctx context.Context, task *asynq.Task) error {
		return tasks.HandleAlarmCheckTask(ctx, task, signalService, alarmService, asyncClient)
	})
	mux.HandleFunc(tasks.TypeSignalWrite, func(ctx context.Context, task *asynq.Task) error {
		return tasks.HandleSignalWriteTask(ctx, task, signalService)
	})
	mux.HandleFunc(tasks.TypeBackofficeCleanSignals, func(ctx context.Context, task *asynq.Task) error {
		return tasks.HandleBackofficeCleanSignalsTask(ctx, task, signalService)
	})
	mux.HandleFunc(tasks.TypeAlarmNotify, func(ctx context.Context, task *asynq.Task) error {
		return tasks.HandleAlarmNotifyTask(ctx, task, alarmService)
	})
}
