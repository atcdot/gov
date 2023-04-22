package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewInstallCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "install",
		Short: "install a version",
		Long:  `Install a version`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please specify a version to install")
				return
			}

			install(args[0])
		},
	}

	return c
}
