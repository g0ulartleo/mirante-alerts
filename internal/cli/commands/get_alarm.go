package commands

import (
	"encoding/json"
	"fmt"

	"github.com/g0ulartleo/mirante-alerts/internal/cli"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
)

type GetAlarmCommand struct{}

func (c *GetAlarmCommand) Name() string {
	return "get-alarm"
}

func (c *GetAlarmCommand) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: ./cli get-alarm <alarm-id>")
	}

	alarmID := args[0]

	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	apiClient := NewAPIClient(cliConfig)

	alarm, err := apiClient.GetAlarm(alarmID)
	if err != nil {
		return fmt.Errorf("failed to get alarm: %w", err)
	}
	jsonData, err := json.MarshalIndent(alarm, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal alarm: %w", err)
	}
	fmt.Printf("%s\n", jsonData)
	return nil
}

func init() {
	c := &GetAlarmCommand{}
	cli.RegisterCommand(c.Name(), c)
}
