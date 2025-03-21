package commands

import (
	"fmt"

	"github.com/g0ulartleo/mirante-alerts/internal/cli"
)

type HelpCommand struct{}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Run(args []string) error {
	commands := []string{}
	for name := range *cli.Registry {
		commands = append(commands, name)
	}
	fmt.Printf("Available commands: %v\n", commands)
	return nil
}

func init() {
	c := &HelpCommand{}
	cli.RegisterCommand(c.Name(), c)
}
