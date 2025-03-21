package commands

import (
	"fmt"

	"github.com/g0ulartleo/mirante-alerts/internal/cli"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
)

type ListAlarmsCommand struct{}

func (c *ListAlarmsCommand) Name() string {
	return "list-alarms"
}

func (c *ListAlarmsCommand) Run(args []string) error {
	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	apiClient := NewAPIClient(cliConfig)
	alarms, err := apiClient.ListAlarms()
	if err != nil {
		return fmt.Errorf("failed to list alarms: %w", err)
	}
	for _, a := range alarms {
		fmt.Printf("%s -> %s\n", a.ID, a.Description)
	}
	return nil
}

func init() {
	c := &ListAlarmsCommand{}
	cli.RegisterCommand(c.Name(), c)
}
