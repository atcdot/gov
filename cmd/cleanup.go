package cmd

import (
	"github.com/spf13/cobra"
)

func NewCleanupCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "cleanup",
		Short: "cleanup alias",
		Long:  `Cleanup alias and use system default go version`,
		Run: func(cmd *cobra.Command, args []string) {
			cleanup()
		},
	}

	return c
}
