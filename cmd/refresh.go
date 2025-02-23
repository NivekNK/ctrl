package cmd

import (
	"ctrl/util"

	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "r",
	Short: "Refresh ctrl system.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fix, err := cmd.Flags().GetBool("fix")
		if err != nil {
			return err
		}

		ctrl, err := util.InitializeInstance()
		if err != nil {
			return err
		}
		defer ctrl.DB.Close()

		config, err := util.LoadConfig()
		if err != nil {
			return err
		}

		data, err := util.LoadData()
		if err != nil {
			return err
		}

		err = data.ForceRefresh(config, ctrl, fix)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	refreshCmd.Flags().BoolP(
		"fix",
		"f",
		false,
		"Fix any problem on the table.",
	)
	rootCmd.AddCommand(refreshCmd)
}
