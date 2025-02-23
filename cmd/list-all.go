package cmd

import (
	"ctrl/util"
	"fmt"

	"github.com/spf13/cobra"
)

var listAllCmd = &cobra.Command{
	Use:   "la",
	Short: "List all apps.",
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
			list, err := ctrl.Query.ListApps(ctrl.Ctx)
			if err != nil {
				return err
			}
			fmt.Print(util.ListAllTable(list))
		} else {
			list, err := ctrl.Query.ListAppsOS(ctrl.Ctx, util.GetOS())
			if err != nil {
				return err
			}
			fmt.Print(util.ListAllTableOS(list))
		}

		return nil
	},
}

func init() {
	listAllCmd.Flags().BoolP(
		"no-os",
		"o",
		false,
		"If used, shows all apps independant of OS.",
	)
	rootCmd.AddCommand(listAllCmd)
}
