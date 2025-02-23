package cmd

import (
	"ctrl/database"
	"ctrl/util"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func updateApp(ctrl *util.Instance, config *util.Config, source string, id string) (bool, error) {
	app, err := ctrl.Query.FindAppBySourceAndId(ctrl.Ctx, database.FindAppBySourceAndIdParams{
		AppSource: source,
		AppID:     id,
		AppOs:     util.GetOS(),
	})
	if err != nil {
		return false, err
	}

	result, err := config.ExecuteCommand(source, util.Versions, id)
	if err != nil {
		return false, err
	}

	versions, err := util.ParseVersions(result)
	if err != nil {
		return false, err
	}

	if !versions.NewVersion.Valid {
		return false, nil
	}

	_, err = config.ExecuteCommand(source, util.Update, id)
	if err != nil {
		return false, err
	}

	err = ctrl.Query.UpdateInstalledApp(ctrl.Ctx, database.UpdateInstalledAppParams{
		AppVersion:   versions.NewVersion,
		AppAvailable: sql.NullString{},
		AppIndex:     app.Index,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

var updateCmd = &cobra.Command{
	Use:   "u [ID / PATH]",
	Short: "Update an app or csv of apps.",
	Args:  cobra.ExactArgs(1),
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

		if isCSV {
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

				updated, err := updateApp(ctrl, config, source, id)
				if err != nil {
					fmt.Println(err.Error())
				}

				if updated {
					fmt.Printf("App %s with source: %s, updated!", id, source)
				}
			}

			return nil
		}

		id := args[0]
		apps, err := ctrl.Query.FindAppById(ctrl.Ctx, database.FindAppByIdParams{
			AppOs: util.GetOS(),
			AppID: id,
		})
		if err != nil {
			return err
		}

		if len(apps) == 0 {
			return fmt.Errorf("app not found")
		}

		if len(apps) == 1 {
			source := apps[0].Source

			updated, err := updateApp(ctrl, config, source, id)
			if err != nil {
				return err
			}

			if updated {
				fmt.Printf("App %s with source: %s, updated!", id, source)
			}
			return nil
		}

		values := []string{}
		for _, app := range apps {
			appSelection := fmt.Sprintf("Id: %s | Source: %s", id, app.Source)
			values = append(values, appSelection)
		}
		values = append(values, "None")

		program := tea.NewProgram(util.NewModel(
			"Select app to install:",
			values,
		))

		model, err := program.Run()
		if err != nil {
			return err
		}

		selected := model.(util.SelectionModel).Selected
		if selected == -1 || selected == len(values)-1 {
			return fmt.Errorf("no app selected to update")
		}

		source := apps[selected].Source
		updated, err := updateApp(ctrl, config, source, id)
		if err != nil {
			return err
		}

		if updated {
			fmt.Printf("App %s with source: %s, updated!", id, source)
		}

		return nil
	},
}

func init() {
	updateCmd.Flags().Bool(
		"csv",
		false,
		"ID is a csv file path.",
	)
	rootCmd.AddCommand(updateCmd)
}
