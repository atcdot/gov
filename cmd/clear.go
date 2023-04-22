package cmd

import (
	"github.com/spf13/cobra"
)

func NewClearCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "clear",
		Short: "clear alias",
		Long:  `Clear alias and use system default go version`,
		Run: func(cmd *cobra.Command, args []string) {
			clear()
		},
	}

	return c
}
