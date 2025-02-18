package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Source struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type Config struct {
	Sources []Source `json:"sources"`
}

func DefaultConfig() *Config {
	config := Config{
		Sources: []Source{
			{Command: "winget", Args: []string{"install"}},
		},
	}
	return &config
}

var ErrLoadingConfig = errors.New("couldnt load ctrl config")

func LoadConfig() (*Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			config := DefaultConfig()

			configFile, err := os.Create(filepath.Join(configPath, "/config.json"))
			if err != nil {
				return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
			}
			defer configFile.Close()

			configJson, err := json.MarshalIndent(config, "", "\t")
			if err != nil {
				return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
			}

			_, err = configFile.Write(configJson)
			if err != nil {
				return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
			}

			fmt.Println("config file created!")
			return config, nil
		} else {
			return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
		}
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("%w :: %s", ErrLoadingConfig, err.Error())
	}

	return &config, nil
}
