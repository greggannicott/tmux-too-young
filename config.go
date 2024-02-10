package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type config struct {
	SearchDirectories []string `yaml:"search_directories"`
}

func getConfig() config {
	homeDir, _ := os.UserHomeDir()
	contents, readFileErr := os.ReadFile(homeDir + "/.tmux-too-young.yaml")
	if readFileErr != nil {
		fmt.Println("Failed to read config file:", readFileErr)
		os.Exit(1)
	}

	var config config
	parseYamlErr := yaml.Unmarshal(contents, &config)
	if parseYamlErr != nil {
		fmt.Println("Failed to parse config file:", parseYamlErr)
	}
	return config
}

func configExists() bool {
	homeDir, _ := os.UserHomeDir()
	_, err := os.Stat(homeDir + "/.tmux-too-young.yaml")
	return err == nil
}

func createConfig(searchDirectoriesString string) {
	searchDirectories := strings.Split(searchDirectoriesString, ",")
	config := config{SearchDirectories: searchDirectories}
	configAsString, _ := yaml.Marshal(config)
	homeDir, _ := os.UserHomeDir()
	os.WriteFile(homeDir+"/.tmux-too-young.yaml", configAsString, 0666)
}
