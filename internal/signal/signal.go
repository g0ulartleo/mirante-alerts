package signal

import "time"

type Signal struct {
	AlertID   string
	Status    Status
	Timestamp time.Time
	Message   string
	Metadata  map[string]interface{}
}

type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusUnknown   Status = "unknown"
)
