package cmd

import (
	"github.com/spf13/cobra"
)

func NewCMD() *cobra.Command {
	c := &cobra.Command{
		Use:   "gov",
		Short: "Go version manager",
		Long:  `Go version manager`,
		// Run: func(cmd *cobra.Command, args []string) {
		// 	listInstalled()
		// },
	}

	c.CompletionOptions.DisableDefaultCmd = true

	return c
}
