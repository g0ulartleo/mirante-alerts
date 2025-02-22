package commands

import "fmt"

type Registry map[string]Command

var commandsRegistry = &Registry{}

func Register(name string, command Command) {
	(*commandsRegistry)[name] = command
}

func Get(name string) (Command, error) {
	command, exists := (*commandsRegistry)[name]
	if !exists {
		return nil, fmt.Errorf("command %s not found", name)
	}
	return command, nil
}
