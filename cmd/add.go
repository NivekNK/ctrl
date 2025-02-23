package cmd

import (
	"ctrl/util"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "a [ID]",
	Short: "Add a new app (not installed).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctrl, err := util.InitializeInstance()
		if err != nil {
			return err
		}
		defer ctrl.DB.Close()

		id := args[0]
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		if len(name) <= 0 {
			name = id
		}

		source, err := cmd.Flags().GetString("source")
		if err != nil {
			return err
		}

		err = ctrl.AddApp(id, name, source)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	addCmd.Flags().StringP(
		"name",
		"n",
		"",
		"Specify a name for the app.",
	)

	addCmd.Flags().StringP(
		"source",
		"s",
		Config.GetDefaultSourceKey(),
		"Specify a source installer for the app.",
	)
	rootCmd.AddCommand(addCmd)
}
