package repo

import (
	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
)

func New() (alarm.AlarmRepository, error) {
	return NewRedisAlarmRepository()
}
