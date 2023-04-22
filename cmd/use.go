package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewUseCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "use",
		Short: "use a version",
		Long:  `Use a version`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please specify a version to use")
				return
			}

			use(args[0])
		},
	}

	return c
}
