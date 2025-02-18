package cmd

import (
	"ctrl/util"
	"fmt"

	"github.com/spf13/cobra"
)

var listNotInstalledCmd = &cobra.Command{
	Use:   "ln",
	Short: "List not installed apps.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctrl, err := InitializeInstance()
		if err != nil {
			return err
		}
		defer ctrl.db.Close()

		no, err := cmd.Flags().GetBool("no-os")
		if err != nil {
			return err
		}

		if no {
			list, err := ctrl.query.FindNotInstalledApps(ctrl.ctx)
			if err != nil {
				return err
			}
			fmt.Print(util.ListNotInstalled(list))
		} else {
			list, err := ctrl.query.FindNotInstalledAppsOS(ctrl.ctx, util.GetOS())
			if err != nil {
				return err
			}
			fmt.Print(util.ListNotInstalledOS(list))
		}
		return nil
	},
}

func init() {
	listNotInstalledCmd.Flags().BoolP(
		"no-os",
		"o",
		false,
		"If used, show not installed apps independant of OS.",
	)
	rootCmd.AddCommand(listNotInstalledCmd)
}
