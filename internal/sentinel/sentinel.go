package sentinel

import (
	"context"
	"time"
)

type Sentinel interface {
	Check(ctx context.Context) (Signal, error)
	Configure(config map[string]interface{}) error
}

// Signal represents a sentinel response
type Signal struct {
	Status    Status
	Timestamp time.Time
	Message   string
	Metadata  map[string]interface{}
}

type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
)

type SentinelConfig struct {
	ID     int
	Name   string
	Type   string
	Config map[string]interface{}
}
