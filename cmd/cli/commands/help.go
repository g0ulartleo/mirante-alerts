package commands

import "fmt"

type HelpCommand struct{}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Run(args []string) error {
	commands := []string{}
	for name := range *commandsRegistry {
		commands = append(commands, name)
	}
	fmt.Printf("Available commands: %v\n", commands)
	return nil
}

func init() {
	c := &HelpCommand{}
	Register(c.Name(), c)
}
