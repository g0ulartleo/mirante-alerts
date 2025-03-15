package web

import (
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	"github.com/g0ulartleo/mirante-alerts/internal/signal"
)

func GetAlarmSignals(signalService *signal.Service, alarmService *alarm.AlarmService) ([]alarm.AlarmSignals, error) {
	alarmsSignals := make([]alarm.AlarmSignals, 0)
	alarms, err := alarmService.GetAlarms()
	if err != nil {
		return nil, err
	}
	for _, a := range alarms {
		signals, err := signalService.GetAlarmLatestSignals(a.ID, 1)
		if err != nil {
			log.Printf("Error fetching signals for alarm %s: %v", a.ID, err)
			signals = []signal.Signal{}
		}
		alarmsSignals = append(alarmsSignals, alarm.AlarmSignals{
			Alarm:   *a,
			Signals: signals,
		})
	}
	return alarmsSignals, nil
}
