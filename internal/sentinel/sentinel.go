package sentinel

import (
	"context"

	"github.com/g0ulartleo/mirante-alerts/internal/signal"
)

type Sentinel interface {
	Check(ctx context.Context, alertID string) (signal.Signal, error)
	Configure(config map[string]interface{}) error
}
