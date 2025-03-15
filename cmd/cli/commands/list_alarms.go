package commands

import (
	"fmt"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
	alarmStores "github.com/g0ulartleo/mirante-alerts/internal/alarm/stores"
)

type ListAlarmsCommand struct{}

func (c *ListAlarmsCommand) Name() string {
	return "list-alarms"
}

func (c *ListAlarmsCommand) Run(args []string) error {
	alarmStore, err := alarmStores.NewAlarmStore()
	if err != nil {
		return fmt.Errorf("failed to initialize alarms: %w", err)
	}
	defer alarmStore.Close()
	alarmService := alarm.NewAlarmService(alarmStore)
	err = alarm.InitAlarms(alarmStore)
	if err != nil {
		return fmt.Errorf("failed to initialize alarms: %w", err)
	}
	alarms, err := alarmService.GetAlarms()
	if err != nil {
		return fmt.Errorf("failed to get alarms: %w", err)
	}
	for _, a := range alarms {
		fmt.Printf("%s -> %s\n", a.ID, a.Description)
	}
	return nil
}

func init() {
	c := &ListAlarmsCommand{}
	Register(c.Name(), c)
}
