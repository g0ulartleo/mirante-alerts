package main

import (
	"fmt"
	"os"

	"github.com/g0ulartleo/mirante-alerts/internal/cli"
	_ "github.com/g0ulartleo/mirante-alerts/internal/cli/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command>")
		os.Exit(1)
	}
	command, err := cli.GetCommand(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := command.Run(os.Args[2:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
