package commands

import (
	"fmt"
	"log"

	"github.com/g0ulartleo/mirante-alerts/internal/cli"
	"github.com/g0ulartleo/mirante-alerts/internal/cli/api_client"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
)

type DeleteAlarmCommand struct{}

func (c *DeleteAlarmCommand) Name() string {
	return "delete-alarm"
}

func (c *DeleteAlarmCommand) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: ./cli %s <alarm-id>", c.Name())
	}
	alarmID := args[0]
	config, err := config.LoadCLIConfig()
	if err != nil {
		return fmt.Errorf("failed to load CLI config: %v", err)
	}
	client := api_client.NewAPIClient(config)
	if err := client.DeleteAlarm(alarmID); err != nil {
		return fmt.Errorf("failed to delete alarm: %v", err)
	}
	log.Printf("Alarm %s deleted", alarmID)
	return nil
}

func init() {
	c := &DeleteAlarmCommand{}
	cli.RegisterCommand(c.Name(), c)
}
