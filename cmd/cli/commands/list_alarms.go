package commands

import (
	"fmt"

	"github.com/g0ulartleo/mirante-alerts/internal/alarm"
)

type ListAlarmsCommand struct{}

func (c *ListAlarmsCommand) Name() string {
	return "list-alarms"
}

func (c *ListAlarmsCommand) Run(args []string) error {
	if err := alarm.InitAlarms(); err != nil {
		return fmt.Errorf("failed to initialize alarms: %w", err)
	}

	for _, a := range alarm.Alarms {
		fmt.Printf("%s -> %s\n", a.ID, a.Description)
	}
	return nil
}

func init() {
	c := &ListAlarmsCommand{}
	Register(c.Name(), c)
}
