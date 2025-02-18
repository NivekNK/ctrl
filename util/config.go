package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Source struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type Config struct {
	Sources      []Source `json:"sources"`
	DelayRefresh int      `json:"delay_refresh"`
}

var defaultConfig = []byte(`{
	"sources": [
		{
			"command": "winget",
			"args": [
				"install"
			]
		}
	],
	"delay_refresh": 168
}`)

func DefaultConfig() *Config {
	var config Config
	err := json.Unmarshal(defaultConfig, &config)
	if err != nil {
		panic(err)
	}

	return &config
}

var ErrLoadingConfig = errors.New("couldnt load ctrl config")

func createConfig(configPath string, config *Config) error {
	configFile, err := os.Create(filepath.Join(configPath, "/config.json"))
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
	}
	defer configFile.Close()

	configJson, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
	}

	_, err = configFile.Write(configJson)
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
	}

	return nil
}

func LoadConfig() (*Config, error) {
	// Get config dirwectory
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	// Set config file path
	configFilePath := filepath.Join(configPath, "/config.json")

	// Check if the file exists
	if _, err := os.Stat(configFilePath); err == nil {
		// Open the config file to read and write
		configFile, err := os.OpenFile(configFilePath, os.O_RDWR, 0644)
		if err != nil {
			return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
		}
		defer configFile.Close()

		// Read all the bytes of the config file
		currentConfigBytes, err := io.ReadAll(configFile)
		if err != nil {
			return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
		}

		// If the config file is empty or {}
		if len(currentConfigBytes) <= 3 {
			// Empty the file
			if err := configFile.Truncate(0); err != nil {
				return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
			}

			if _, err := configFile.Seek(0, 0); err != nil {
				return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
			}

			// Write default config
			if _, err = configFile.Write(defaultConfig); err != nil {
				return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
			}

			// Get default config to config struct
			var config Config
			if err := json.Unmarshal(defaultConfig, &config); err != nil {
				return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
			}

			fmt.Println("config file updated!")
			return &config, nil
		}

		// Get the current config into a map
		var currentConfigMap map[string]interface{}
		if err := json.Unmarshal(currentConfigBytes, &currentConfigMap); err != nil {
			return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
		}

		// Get the default config into a map
		var defaultConfigMap map[string]interface{}
		if err := json.Unmarshal(defaultConfig, &defaultConfigMap); err != nil {
			return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
		}

		// Iterate over the default config to check if
		// the current config doesn`t have a default key;
		// updated if it doesn`t exists
		updatedConfig := false
		for key := range defaultConfigMap {
			if _, exists := currentConfigMap[key]; !exists {
				currentConfigMap[key] = defaultConfigMap[key]
				updatedConfig = true
			}
		}

		// If the current config was not updated
		// returns the current config
		if !updatedConfig {
			var config Config
			if err := json.Unmarshal(currentConfigBytes, &config); err != nil {
				return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
			}
			return &config, nil
		}

		// Tranforms updated config into bytes
		newConfigBytes, err := json.MarshalIndent(currentConfigMap, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
		}

		// Empty config file
		if err := configFile.Truncate(0); err != nil {
			return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
		}

		if _, err := configFile.Seek(0, 0); err != nil {
			return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
		}

		// Write updated config to file
		if _, err = configFile.Write(newConfigBytes); err != nil {
			return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
		}

		// Get the new config to config struct
		var config Config
		err = json.Unmarshal(newConfigBytes, &config)
		if err != nil {
			return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
		}

		fmt.Println("config file updated!")
		return &config, nil
	} else if !os.IsNotExist(err) {
		config := DefaultConfig()

		err := createConfig(configPath, config)
		if err != nil {
			return nil, err
		}

		fmt.Println("config file created!")
		return config, nil
	} else {
		return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
	}
}
