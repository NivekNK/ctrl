package cmd

import (
	"ctrl/util"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "a [SOURCE(?)] [ID / PATH]",
	Short: "Add a new app or csv of apps (not installed).",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctrl, err := util.InitializeInstance()
		if err != nil {
			return err
		}
		defer ctrl.DB.Close()

		isCSV, err := cmd.Flags().GetBool("csv")
		if err != nil {
			return err
		}

		if len(args) == 1 && isCSV {
			path := args[0]

			file, err := os.Open(path)
			if err != nil {
				return err
			}

			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				return err
			}

			for _, row := range records {
				source := row[0]
				id := row[1]
				name := row[2]

				err = ctrl.AddApp(id, name, source)
				if err != nil {
					fmt.Println(err.Error())
				}
			}

			return nil
		}

		var source string
		var id string

		if len(args) == 1 {
			config, err := util.LoadConfig()
			if err != nil {
				return nil
			}

			source = config.GetDefaultSourceKey()
			id = args[0]
		} else {
			source = args[0]
			id = args[1]
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		if len(name) <= 0 {
			name = id
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
	addCmd.Flags().Bool(
		"csv",
		false,
		"ID is a csv file path.",
	)
	rootCmd.AddCommand(addCmd)
}
