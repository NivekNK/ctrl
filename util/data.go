package util

import (
	"ctrl/database"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Data struct {
	LastRefresh string `json:"last_refresh"`
}

func DefaultData() *Data {
	data := Data{
		LastRefresh: time.Now().Format(time.RFC3339),
	}
	return &data
}

var ErrLoadingData = errors.New("couldnt load ctrl data")

func createData(dataFilePath string, data *Data) error {
	dataFile, err := os.Create(dataFilePath)
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrLoadingData, err.Error())
	}
	defer dataFile.Close()

	dataJson, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrLoadingData, err.Error())
	}

	_, err = dataFile.Write(dataJson)
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrLoadingData, err.Error())
	}

	return nil
}

func LoadData() (*Data, error) {
	dataPath, err := DataPath()
	if err != nil {
		return nil, err
	}

	dataFilePath := filepath.Join(dataPath, "/data.json")
	if _, err := os.Stat(dataFilePath); err == nil {
		dataBytes, err := os.ReadFile(dataFilePath)
		if err != nil {
			return nil, fmt.Errorf("%w :: %s", ErrLoadingData, err.Error())
		}

		var data Data
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return nil, fmt.Errorf("%w :: %s", ErrLoadingData, err.Error())
		}

		return &data, nil
	} else if os.IsNotExist(err) {
		data := DefaultData()

		err := createData(dataFilePath, data)
		if err != nil {
			return nil, err
		}

		fmt.Println("data file created!")
		return data, nil
	} else {
		return nil, fmt.Errorf("%w :: %s", ErrLoadingData, err.Error())
	}
}

var ErrRefreshError = errors.New("refresh problem")

func (data *Data) ForceRefresh(config *Config, instance *Instance, fix bool) error {
	apps, err := instance.Query.ListAppsOS(instance.Ctx, GetOS())
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrRefreshError, err.Error())
	}

	for _, app := range apps {
		result, err := config.ExecuteCommand(app.Source, Versions, app.ID)
		if err != nil {
			return err
		}

		versions, err := ParseVersions(result)
		if err != nil {
			return err
		}

		if app.Installed && !fix {
			if !versions.NewVersion.Valid {
				continue
			}

			err = instance.Query.UpdateAvailable(instance.Ctx, database.UpdateAvailableParams{
				AppAvailable: versions.NewVersion,
				AppIndex:     app.Index,
			})
			if err != nil {
				return fmt.Errorf("%w :: %s", ErrRefreshError, err.Error())
			}
		} else {
			err = instance.Query.UpdateInstalledApp(instance.Ctx, database.UpdateInstalledAppParams{
				AppVersion:   sql.NullString{String: versions.Version, Valid: true},
				AppAvailable: versions.NewVersion,
				AppIndex:     app.Index,
			})
			if err != nil {
				return fmt.Errorf("%w :: %s", ErrRefreshError, err.Error())
			}
		}
	}

	fmt.Println("Refresh successful!")
	return nil
}

func (data *Data) AdviceRefresh(config *Config) error {
	if config == nil {
		return fmt.Errorf("%w :: %s", ErrRefreshError, "invalid config")
	}

	lastRefresh, err := time.Parse(time.RFC3339, data.LastRefresh)
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrRefreshError, err.Error())
	}

	duration := time.Since(lastRefresh)
	if config.Refresh.AdviceDelay == -1 || duration.Hours() < float64(config.Refresh.AdviceDelay) {
		return nil
	}

	if config.Refresh.ForceRefreshDelay != -1 && duration.Hours() >= float64(config.Refresh.ForceRefreshDelay) {
		ctrl, err := InitializeInstance()
		if err != nil {
			return fmt.Errorf("%w :: %s", ErrRefreshError, err.Error())
		}
		defer ctrl.DB.Close()

		fmt.Println("The configured time has already passed, forcing refresh...")
		data.ForceRefresh(config, ctrl, true)
		return nil
	}

	fmt.Println("The configured time has already passed, you should refresh!")

	return nil
}
