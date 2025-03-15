package tasks

import (
	"context"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	alarmTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/alarm"
	backofficeTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/backoffice"
	sigTasks "github.com/g0ulartleo/mirante-alerts/internal/worker/tasks/signal"
	"github.com/hibiken/asynq"
)

func Register(mux *asynq.ServeMux, signalService *signal.Service, alarmService *alarm.AlarmService, asyncClient *asynq.Client) {
	mux.HandleFunc(alarmTasks.TypeCheckAlarm, func(ctx context.Context, task *asynq.Task) error {
		return alarmTasks.HandleCheckAlarmTask(ctx, task, signalService, alarmService, asyncClient)
	})
	mux.HandleFunc(sigTasks.TypeSignalWrite, func(ctx context.Context, task *asynq.Task) error {
		return sigTasks.HandleSignalWriteTask(ctx, task, signalService)
	})
	mux.HandleFunc(backofficeTasks.TypeCleanSignals, func(ctx context.Context, task *asynq.Task) error {
		return backofficeTasks.HandleCleanSignalsTask(ctx, task, signalService)
	})
	mux.HandleFunc(alarmTasks.TypeNotify, func(ctx context.Context, task *asynq.Task) error {
		return alarmTasks.HandleNotifyTask(ctx, task, alarmService)
	})
}
