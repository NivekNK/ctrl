package cmd

import (
	"ctrl/util"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ctrl",
	Short: "Cross-Platform app management.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := util.LoadConfig()
		if err != nil {
			return err
		}

		data, err := util.LoadData()
		if err != nil {
			return err
		}
		data.AdviceRefresh(config)

		return cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
