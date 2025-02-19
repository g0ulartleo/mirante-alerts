package alarm

import "github.com/g0ulartleo/mirante-alerts/internal/signal"

type Alarm struct {
	ID          string
	Name        string
	Description string
	Path        []string
	Type        string
	Config      map[string]interface{}
	Cron        string
	Interval    string
}

type AlarmSignals struct {
	Alarm   Alarm
	Signals []signal.Signal
}
