package main

import (
	"os"

	"github.com/atcdot/gov/cmd"
)

func main() {
	command := cmd.NewCMD()

	command.AddCommand(cmd.NewInitCmd())
	command.AddCommand(cmd.NewGUICmd())
	command.AddCommand(cmd.NewListCmd())
	command.AddCommand(cmd.NewInstallCmd())
	command.AddCommand(cmd.NewRemoveCmd())
	command.AddCommand(cmd.NewUseCmd())
	command.AddCommand(cmd.NewCleanupCmd())

	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
