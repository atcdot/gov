package cmd

import (
	"github.com/spf13/cobra"

	"gov/internal/gui"
)

func NewGUICmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "gui",
		Short: "GUI mode",
		Long:  `Start in GUI mode`,
		Run: func(cmd *cobra.Command, args []string) {
			gui.RunGUI()
		},
	}

	return c
}
