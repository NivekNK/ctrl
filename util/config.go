package util

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Source struct {
	Versions  []string `json:"versions"`
	Install   []string `json:"install"`
	Update    []string `json:"update"`
	Uninstall []string `json:"uninstall"`
}

type Refresh struct {
	AdviceDelay       int `json:"advice_delay"`
	ForceRefreshDelay int `json:"force_refresh_delay"`
}

type Config struct {
	Sources map[string]Source `json:"sources"`
	Refresh Refresh           `json:"refresh"`
}

var defaultConfig = []byte(`{
	"sources": {
		"winget": {
			"versions": [
				"pwsh",
				"-c",
				"$wingetApp = \"APP_ID\";",
				"$wingetList = winget list $wingetApp;",
				"if ($wingetList -match \"No installed package found matching input criteria.\") { \"\" } else {",
				"    $array = $wingetList[$wingetList.Length - 1] -split \"\\s+\";",
				"    [array]::Reverse($array);",
				"    $count = ($wingetList[$wingetList.Length - 3] -split \"\\s+\").Length;",
				"    if ($count -ge 5) { $array[2] + \",\" + $array[1] } else { $array[1] + \",_\" }",
				"}"
			],
			"install": [
				"winget",
				"install",
				"APP_ID"
			],
			"update": [
				"winget",
				"upgrade",
				"--id",
				"APP_ID"
			],
			"uninstall": [
				"winget",
				"uninstall",
				"APP_ID"
			]
		}
	},
	"refresh": {
		"advice_delay": 168,
		"force_refresh_delay": 672
	}
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

func createConfig(configFilePath string, config *Config) error {
	configFile, err := os.Create(configFilePath)
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
	} else if os.IsNotExist(err) {
		config := DefaultConfig()

		err := createConfig(configFilePath, config)
		if err != nil {
			return nil, err
		}

		fmt.Println("config file created!")
		return config, nil
	} else {
		return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
	}
}

func (config *Config) GetDefaultSourceKey() string {
	var defaultKey string = ""
	for key := range config.Sources {
		defaultKey = key
		break
	}
	return defaultKey
}

type CmdType int

const (
	Versions CmdType = iota
	Install
	Update
	Uninstall
)

func (cmdType CmdType) String() string {
	switch cmdType {
	case Versions:
		return "Versions"
	case Install:
		return "Install"
	case Update:
		return "Update"
	case Uninstall:
		return "Uninstall"
	default:
		return "Invalid"
	}
}

func (config *Config) ExecuteCommand(source string, cmdType CmdType, appId string) (string, error) {
	selected, exists := config.Sources[source]
	if !exists {
		return "", fmt.Errorf("%w :: not configured app source (%s) in app: %s", ErrRefreshError, source, appId)
	}

	var commandName string
	var commands []string

	switch cmdType {
	case Versions:
		commandName = selected.Versions[0]
		commands = append([]string{}, selected.Versions[1:]...)
	case Install:
		commandName = selected.Install[0]
		commands = append([]string{}, selected.Install[1:]...)
	case Update:
		commandName = selected.Update[0]
		commands = append([]string{}, selected.Update[1:]...)
	case Uninstall:
		commandName = selected.Uninstall[0]
		commands = append([]string{}, selected.Uninstall[1:]...)
	default:
		return "", fmt.Errorf("%w :: wrong cmd type", ErrRefreshError)
	}

	found := false
	for i := len(commands) - 1; i >= 0; i-- {
		fmt.Println(commands[i])
		if strings.Contains(commands[i], "APP_ID") {
			commands[i] = strings.Replace(commands[i], "APP_ID", appId, 1)
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("%w :: cmd type doesn`t have app_id: %s", ErrRefreshError, cmdType.String())
	}

	command := exec.Command(commandName, commands...)
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("%w :: %s", ErrRefreshError, err.Error())
	}

	return string(output), nil
}

type ParsedVersion struct {
	Version    string
	NewVersion sql.NullString
}

var ErrParseVersion = errors.New("problem parsing versions command")

func ParseVersions(versions string) (*ParsedVersion, error) {
	if len(versions) == 0 {
		return nil, fmt.Errorf("%w :: couldn`t find app", ErrParseVersion)
	}

	versionParts := strings.Split(versions, ",")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("%w :: versions command with incorrect format", ErrParseVersion)
	}

	newVersion := versionParts[1]
	var newVersionValue sql.NullString
	if newVersion != "_" {
		newVersionValue.String = newVersion
		newVersionValue.Valid = true
	}

	parsedVersion := ParsedVersion{
		Version:    versionParts[0],
		NewVersion: newVersionValue,
	}
	return &parsedVersion, nil
}
