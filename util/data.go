package util

import (
	"ctrl/database"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	} else if !os.IsNotExist(err) {
		data := DefaultData()

		err := createData(dataFilePath, data)
		if err != nil {
			return nil, err
		}

		return data, nil
	} else {
		return nil, fmt.Errorf("%w :: %s", ErrLoadingData, err.Error())
	}
}

var ErrRefreshError = errors.New("refresh problem")

func (data *Data) ForceRefresh(config *Config, instance *Instance) error {
	apps, err := instance.Query.ListAppsOS(instance.Ctx, GetOS())
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrRefreshError, err.Error())
	}

	for _, app := range apps {
		versions, err := config.ExecuteCommand(app.Source, Versions, app.ID)
		if err != nil {
			return err
		}

		versionParts := strings.Split(versions, ",")
		if len(versionParts) != 2 {
			return fmt.Errorf("%w :: versions command for source (%s) incorrect format", ErrRefreshError, app.Source)
		}

		version := versionParts[0]
		newVersion := versionParts[1]

		if app.Installed {
			if newVersion == "_" {
				return nil
			}

			_, err = config.ExecuteCommand(app.Source, Update, app.ID)
			if err != nil {
				return fmt.Errorf("%w :: %s", ErrRefreshError, err.Error())
			}

			err = instance.Query.UpdateAvailable(instance.Ctx, database.UpdateAvailableParams{
				AppAvailable: sql.NullString{String: newVersion},
				AppIndex:     app.Index,
			})
			if err != nil {
				return fmt.Errorf("%w :: %s", ErrRefreshError, err.Error())
			}

			return nil
		} else {
			// TODO: FIX THIS SHIT
			if len(version) == 0 {
				return fmt.Errorf("%w :: couldnt find app from source (%s) with id: %s", ErrRefreshError, app.Source, app.ID)
			}

			err = instance.Query.UpdateInstalledApp(instance.Ctx, database.UpdateInstalledAppParams{
				AppVersion: sql.NullString{String: version},
				AppIndex:   app.Index,
			})
			if err != nil {
				return fmt.Errorf("%w :: %s", ErrRefreshError, err.Error())
			}

			// Revisar si una app agregada esta instalada, la version y actualizar la avaliable
			return nil
		}
	}

	return nil
}

func (data *Data) Refresh(config *Config, instance *Instance) error {
	if config == nil {
		return fmt.Errorf("%w :: %s", ErrRefreshError, "invalid config")
	}

	lastRefresh, err := time.Parse(time.RFC3339, data.LastRefresh)
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrRefreshError, err.Error())
	}

	duration := time.Since(lastRefresh)
	if duration.Hours() < float64(config.DelayRefresh) {
		return nil
	}

	return data.ForceRefresh(config, instance)
}
