package cmd

import (
	"ctrl/util"
	"fmt"

	"github.com/spf13/cobra"
)

var listInstalledCmd = &cobra.Command{
	Use:   "li",
	Short: "List installed apps.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctrl, err := util.InitializeInstance()
		if err != nil {
			return err
		}
		defer ctrl.DB.Close()

		no, err := cmd.Flags().GetBool("no-os")
		if err != nil {
			return err
		}

		if no {
			list, err := ctrl.Query.FindInstalledApps(ctrl.Ctx)
			if err != nil {
				return err
			}
			fmt.Print(util.ListInstalled(list))
		} else {
			list, err := ctrl.Query.FindInstalledAppsOS(ctrl.Ctx, util.GetOS())
			if err != nil {
				return err
			}
			fmt.Print(util.ListInstalledOS(list))
		}
		return nil
	},
}

func init() {
	listInstalledCmd.Flags().BoolP(
		"no-os",
		"o",
		false,
		"If used, show installed apps independant of OS.",
	)
	rootCmd.AddCommand(listInstalledCmd)
}
