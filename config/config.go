package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	SearchDirectories []string `yaml:"search_directories"`
}

func GetConfig() Configuration {
	homeDir, _ := os.UserHomeDir()
	contents, readFileErr := os.ReadFile(homeDir + "/.tmux-too-young.yaml")
	if readFileErr != nil {
		fmt.Println("Failed to read config file:", readFileErr)
		os.Exit(1)
	}

	var config Configuration
	parseYamlErr := yaml.Unmarshal(contents, &config)
	if parseYamlErr != nil {
		fmt.Println("Failed to parse config file:", parseYamlErr)
	}
	return config
}

func ConfigExists() bool {
	homeDir, _ := os.UserHomeDir()
	_, err := os.Stat(homeDir + "/.tmux-too-young.yaml")
	return err == nil
}

func CreateConfig(searchDirectoriesString []string) {
	config := Configuration{SearchDirectories: searchDirectoriesString}
	configAsString, _ := yaml.Marshal(config)
	homeDir, _ := os.UserHomeDir()
	os.WriteFile(homeDir+"/.tmux-too-young.yaml", configAsString, 0666)
}
