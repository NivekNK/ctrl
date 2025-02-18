package cmd

import (
	"ctrl/util"
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "a ID",
	Short: "Add a new app (not installed).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctrl, err := InitializeInstance()
		if err != nil {
			return err
		}
		defer ctrl.db.Close()

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

	var command string = ""

	config, err := util.LoadConfig()
	if err != nil {
		fmt.Println(err.Error())
	} else if len(config.Sources) > 0 {
		command = config.Sources[0].Command
	}

	addCmd.Flags().StringP(
		"source",
		"s",
		command,
		"Specify a source installer for the app.",
	)
	rootCmd.AddCommand(addCmd)
}
