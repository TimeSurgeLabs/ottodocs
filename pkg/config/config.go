package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
)

// Config represents the configuration file
type Config struct {
	APIKey string `json:"api_key"`
	Model  string `json:"model"`
}

// also returns the path to the config file
func createIfNotExists() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	homeDir := currentUser.HomeDir
	configDir := filepath.Join(homeDir, ".ottodocs")
	configPath := filepath.Join(configDir, "config.json")

	// check if the config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// create the config directory
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			return "", err
		}
	}

	// check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// create the config file
		file, err := os.Create(configPath)
		if err != nil {
			return "", err
		}
		// add an empty config to the file
		blankConfig := Config{
			APIKey: "",
			Model:  "gpt-3.5-turbo",
		}
		err = json.NewEncoder(file).Encode(blankConfig)
		if err != nil {
			return "", err
		}
		defer file.Close()
	}

	return configPath, nil
}

// Load loads the configuration file
func Load() (*Config, error) {
	// load the config file at path ~/.ottodocs/config.json
	// if the file or path does not exist, create it
	configPath, err := createIfNotExists()
	if err != nil {
		return nil, err
	}

	// open the config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	// decode the config file
	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Save saves the configuration file
func (c *Config) Save() error {
	configPath, err := createIfNotExists()
	if err != nil {
		return err
	}

	// open the config file
	file, err := os.OpenFile(configPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	return json.NewEncoder(file).Encode(c)
}

// Identical save function except that it takes the config object as its argument
func Save(c *Config) error {
	configPath, err := createIfNotExists()
	if err != nil {
		return err
	}

	// open the config file
	file, err := os.OpenFile(configPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	return json.NewEncoder(file).Encode(c)
}
