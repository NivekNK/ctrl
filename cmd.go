package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ctrl",
	Short: "Cross-Platform app management.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var addCmd = &cobra.Command{
	Use:   "add ID",
	Short: "Add a new app (not installed).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctrl, err := initDB()
		if err != nil {
			return err
		}
		defer ctrl.instance.Close()

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

		err = addApp(ctrl.ctx, ctrl.db, id, name, source)
		if err != nil {
			return err
		}

		return nil
	},
}

var listAllCmd = &cobra.Command{
	Use:   "la",
	Short: "List all apps.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctrl, err := initDB()
		if err != nil {
			return err
		}
		defer ctrl.instance.Close()

		no, err := cmd.Flags().GetBool("no-os")
		if err != nil {
			return err
		}

		if no {
			list, err := ctrl.db.ListApps(ctrl.ctx)
			if err != nil {
				return err
			}
			fmt.Print(listAllTable(list))
		} else {
			list, err := ctrl.db.ListAppsOS(ctrl.ctx, getOS())
			if err != nil {
				return err
			}
			fmt.Print(listAllTableOS(list))
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
		"winget",
		"Specify a source installer for the app.",
	)
	rootCmd.AddCommand(addCmd)

	listAllCmd.Flags().BoolP(
		"no-os",
		"o",
		false,
		"If used, shows all apps independant of OS.",
	)
	rootCmd.AddCommand(listAllCmd)
}
