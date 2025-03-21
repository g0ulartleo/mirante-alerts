package commands

import (
	"fmt"

	"github.com/g0ulartleo/mirante-alerts/internal/cli"
	"github.com/g0ulartleo/mirante-alerts/internal/config"
)

type ConfigCommand struct{}

func (c *ConfigCommand) Name() string {
	return "config"
}

func (c *ConfigCommand) Run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: ./cli config <api_host> <api_key>")
	}

	apiHost := args[0]
	apiKey := args[1]

	cliConfig := &config.CLIConfig{
		APIHost: apiHost,
		APIKey:  apiKey,
	}

	return config.SaveCLIConfig(cliConfig)
}

func init() {
	c := &ConfigCommand{}
	cli.RegisterCommand(c.Name(), c)
}
