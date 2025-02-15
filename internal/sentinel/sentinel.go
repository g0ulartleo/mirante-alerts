package sentinel

import (
	"context"

	"github.com/g0ulartleo/mirante-alerts/internal/signal"
)

type Sentinel interface {
	Check(ctx context.Context) (signal.Signal, error)
	Configure(config map[string]interface{}) error
}

type SentinelConfig struct {
	ID       string
	Name     string
	Path     []string
	Type     string
	Config   map[string]interface{}
	Cron     string
	Interval string
}

type SentinelConfigData struct {
	Config  SentinelConfig
	Signals []signal.Signal
}
