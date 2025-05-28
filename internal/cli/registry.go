package cli

import "fmt"

type Command interface {
	Name() string
	Description() string
	Usage() string
	Run(args []string) error
}

type CommandsRegistry map[string]Command

var Registry = &CommandsRegistry{}

func RegisterCommand(name string, command Command) {
	(*Registry)[name] = command
}

func GetCommand(name string) (Command, error) {
	command, exists := (*Registry)[name]
	if !exists {
		return nil, fmt.Errorf("command %s not found", name)
	}
	return command, nil
}
