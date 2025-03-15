package stores

import (
	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
)

func NewAlarmStore() (alarm.AlarmRepository, error) {
	return NewRedisAlarmRepository()
}
