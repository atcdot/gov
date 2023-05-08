package cmd

import (
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "init gov",
		Long:  `Creates config`,
		Run: func(cmd *cobra.Command, args []string) {
			initialise()
		},
	}

	return c
}
