package notification

import (
	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
)

type Notification interface {
	Build(alarmConfig *alarm.Alarm, sig signal.Signal) error
	Send() error
}
