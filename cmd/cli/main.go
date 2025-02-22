package main

import (
	"fmt"
	"os"

	"github.com/g0ulartleo/mirante-alerts/cmd/cli/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command>")
		os.Exit(1)
	}
	command, err := commands.Get(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := command.Run(os.Args[2:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
