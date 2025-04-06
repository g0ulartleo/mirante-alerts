package worker

import (
	"context"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/g0ulartleo/mirante-alerts/internal/worker/tasks"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

func RegisterTasks(mux *asynq.ServeMux, sentinelFactory *sentinel.SentinelFactory, signalService *signal.Service, alarmService *alarm.AlarmService, asyncClient *asynq.Client, redisClient *redis.Client) {
	mux.HandleFunc(tasks.TypeAlarmCheck, func(ctx context.Context, task *asynq.Task) error {
		return tasks.HandleAlarmCheckTask(ctx, task, sentinelFactory, signalService, alarmService, asyncClient)
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
	mux.HandleFunc(tasks.TypeDashboardNotify, func(ctx context.Context, task *asynq.Task) error {
		return tasks.HandleDashboardNotifyTask(ctx, task, signalService, alarmService, redisClient)
	})
}
