package cmd

import (
	"ctrl/database"
	"ctrl/util"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func installApp(ctrl *util.Instance, config *util.Config, source string, id string, name string) error {
	result, _ := config.ExecuteCommand(source, util.Versions, id)
	versions, parseErr := util.ParseVersions(result)

	app, findAppErr := ctrl.Query.FindAppBySourceAndId(ctrl.Ctx, database.FindAppBySourceAndIdParams{
		AppSource: source,
		AppID:     id,
		AppOs:     util.GetOS(),
	})

	if parseErr == nil && findAppErr != nil {
		err := ctrl.InstallApp(id, name, source, versions.Version, versions.NewVersion)
		if err != nil {
			return fmt.Errorf("failed to add installed app :: %w", err)
		}

		fmt.Printf("App %s with source: %s, already installed. Added to the database.\n", id, source)
		return nil
	} else if parseErr == nil && findAppErr == nil {
		if versions.NewVersion.Valid {
			err := ctrl.Query.UpdateAvailable(ctrl.Ctx, database.UpdateAvailableParams{
				AppAvailable: versions.NewVersion,
				AppIndex:     app.Index,
			})
			if err != nil {
				return fmt.Errorf("failed to updated added app :: %w", err)
			}

			fmt.Printf("App %s with source: %s, already installed. New version available.\n", id, source)
			return nil
		}

		fmt.Printf("App %s with source: %s, already installed.\n", id, source)
		return nil
	} else if parseErr != nil && findAppErr == nil {
		fmt.Printf("App %s with source: %s, already installed. Error: %s\n", id, source, parseErr.Error())
		return nil
	}

	_, err := config.ExecuteCommand(source, util.Install, id)
	if err != nil {
		return fmt.Errorf("failed to install app :: %w", err)
	}

	result, err = config.ExecuteCommand(source, util.Versions, id)
	if err != nil {
		return fmt.Errorf("failed to add installed app :: %w", err)
	}

	versions, err = util.ParseVersions(result)
	if err != nil {
		return err
	}

	err = ctrl.InstallApp(id, name, source, versions.Version, versions.NewVersion)
	if err != nil {
		return fmt.Errorf("failed to add installed app :: %w", err)
	}

	fmt.Printf("%s installed!\n", id)
	return nil
}

var installCmd = &cobra.Command{
	Use:   "i [SOURCE(?)] [ID / PATH]",
	Short: "Add and install an app or csv of apps.",
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

				err = installApp(ctrl, config, source, id, name)
				if err != nil {
					fmt.Println(err.Error())
				}
			}

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

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		if len(name) == 0 {
			name = id
		}

		return installApp(ctrl, config, source, id, name)
	},
}

func init() {
	installCmd.Flags().StringP(
		"name",
		"n",
		"",
		"Specify a name for the app.",
	)
	installCmd.Flags().Bool(
		"csv",
		false,
		"ID is a csv file path.",
	)
	installCmd.MarkFlagsMutuallyExclusive("name", "csv")
	rootCmd.AddCommand(installCmd)
}
