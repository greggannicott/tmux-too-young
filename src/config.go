package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	SearchDirectories []string `yaml:"search_directories"`
}

func getConfig() config {
	homeDir, _ := os.UserHomeDir()
	contents, readFileErr := os.ReadFile(homeDir + "/.ts.yaml")
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
