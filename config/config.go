package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

const (
	configRelPath        = "pact/config.toml"
	defaultScriptsRelDir = "pact/scripts"
)

type Config struct {
	ScriptsDir string `toml:"scripts_dir"`
}

func DefaultConfig() *Config {
	scriptsDir := xdg.DataHome + "/" + defaultScriptsRelDir
	return &Config{
		ScriptsDir: scriptsDir,
	}
}

func NewConfig() (*Config, error) {
	// run setup seen in example folder
	if os.Getenv("EXAMPLE") == "1" {
		return ExampleConfig()
	}

	filePath, err := xdg.SearchConfigFile(configRelPath)
	if err != nil {
		fmt.Println("no config was found")
		newFilePath, err := createConfigFile()
		if err != nil {
			fmt.Println("error creating config")
			return nil, err
		}
		filePath = newFilePath
	}

	var config Config
	_, err = toml.DecodeFile(filePath, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func createConfigFile() (string, error) {
	filePath, err := xdg.ConfigFile(configRelPath)
	if err != nil {
		return "", err
	}
	fmt.Println("creating default config file at ", filePath)
	cfg := DefaultConfig()

	tomlContents, err := toml.Marshal(cfg)
	if err != nil {
		return "", err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	_, err = f.Write(tomlContents)
	if err != nil {
		return "", err
	}

	fmt.Println("config file is created")
	return filePath, nil
}

func ExampleConfig() (*Config, error) {
	// running from the root of the project
	filePath := "./example/config.toml"

	var config Config
	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
