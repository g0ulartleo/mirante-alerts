package commands

import (
	"fmt"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/cli"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
)

type CheckAlarmCommand struct{}

func (c *CheckAlarmCommand) Name() string {
	return "check-alarm"
}

func (c *CheckAlarmCommand) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: ./cli %s <alarm-id>", c.Name())
	}

	alarmID := args[0]
	config, err := config.LoadCLIConfig()
	if err != nil {
		return fmt.Errorf("failed to load CLI config: %v", err)
	}
	client := NewAPIClient(config)
	if err := client.CheckAlarm(alarmID); err != nil {
		return fmt.Errorf("failed to check alarm: %v", err)
	}

	log.Printf("Task enqueued")
	return nil
}

func init() {
	c := &CheckAlarmCommand{}
	cli.RegisterCommand(c.Name(), c)
}
