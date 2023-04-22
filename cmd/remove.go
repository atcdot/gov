package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRemoveCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "remove",
		Short: "remove a version",
		Long:  `Remove a version`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please specify a version to remove")
				return
			}

			remove(args[0])
		},
	}

	return c
}
