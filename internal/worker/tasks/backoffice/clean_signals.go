package backoffice

import (
	"context"

	"github.com/g0ulartleo/mirante-alerts/internal/signal"
	"github.com/hibiken/asynq"
)

const (
	TypeCleanSignals = "backoffice:clean-signals"
)

func NewCleanSignalsTask() (*asynq.Task, error) {
	return asynq.NewTask(TypeCleanSignals, nil, asynq.MaxRetry(3)), nil
}

func HandleCleanSignalsTask(ctx context.Context, t *asynq.Task, signalService *signal.Service) error {
	return signalService.CleanOldSignals()
}
