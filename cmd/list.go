package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "list installed and available versions",
		Long:  `List installed and available versions`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Installed versions:")
			listInstalled()

			fmt.Println("Available versions:")
			listAvailable()
		},
	}

	return c
}
