package main

import (
	"os"

	"gov/cmd"
)

func main() {
	command := cmd.NewCMD()

	command.AddCommand(cmd.NewGUICmd())
	command.AddCommand(cmd.NewListCmd())
	command.AddCommand(cmd.NewInstallCmd())
	command.AddCommand(cmd.NewRemoveCmd())
	command.AddCommand(cmd.NewUseCmd())
	command.AddCommand(cmd.NewClearCmd())

	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
