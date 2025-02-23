package cmd

import (
	"ctrl/util"
	"os"

	"github.com/spf13/cobra"
)

var Config *util.Config = nil

var rootCmd = &cobra.Command{
	Use:   "ctrl",
	Short: "Cross-Platform app management.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
