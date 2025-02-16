package alert

import "github.com/g0ulartleo/mirante-alerts/internal/signal"

type Alert struct {
	ID       string
	Name     string
	Path     []string
	Type     string
	Config   map[string]interface{}
	Cron     string
	Interval string
}

type AlertSignals struct {
	Alert   Alert
	Signals []signal.Signal
}
