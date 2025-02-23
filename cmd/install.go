package cmd

import (
	"ctrl/database"
	"ctrl/util"
	"fmt"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "i [SOURCE(?)] [ID]",
	Short: "Add and install an app.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctrl, err := util.InitializeInstance()
		if err != nil {
			return err
		}
		defer ctrl.DB.Close()

		config, err := util.LoadConfig()
		if err != nil {
			return nil
		}

		var source string
		var id string

		if len(args) == 1 {
			source = config.GetDefaultSourceKey()
			id = args[0]
		} else {
			source = args[0]
			id = args[1]
		}

		app, err := ctrl.Query.FindAppBySourceAndId(ctrl.Ctx, database.FindAppBySourceAndIdParams{
			AppSource: source,
			AppID:     id,
			AppOs:     util.GetOS(),
		})

		if err == nil {
			if app.Installed {
				if app.Available.Valid {
					fmt.Println("app already installed, update available!")
				} else {
					fmt.Println("app already installed!")
				}
				return nil
			}
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		_, err = config.ExecuteCommand(source, util.Install, id)
		if err != nil {
			return err
		}

		if len(name) == 0 {
			name = id
		}

		result, err := config.ExecuteCommand(source, util.Versions, id)
		if err != nil {
			return err
		}

		versions, err := util.ParseVersions(result)
		if err != nil {
			return err
		}

		err = ctrl.InstallApp(id, name, source, versions.Version)
		if err != nil {
			return err
		}

		fmt.Printf("%s installed!\n", id)

		return nil
	},
}

func init() {
	installCmd.Flags().StringP(
		"name",
		"n",
		"",
		"Specify a name for the app.",
	)
	rootCmd.AddCommand(installCmd)
}
