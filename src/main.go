package main

import (
	"log/slog"
	"os"
	"strings"
)

type worktreeDetails struct {
	worktree string
	branch   string
}

var projects []project
var logger *slog.Logger

// To run in terminal: go run tmux-too-young
// Config file: ~/.tmux-too-young.yaml
func main() {
	setupLogging()
	initialSearchTerm := getInitialSearchTerm()
	config := getConfig()
	scanProjectDirectories(config)
	selectedProject := getSelectionFromUser(initialSearchTerm)
	launchProject(selectedProject)
}

func getInitialSearchTerm() string {
	st := strings.Join(os.Args[1:], " ")
	if st != "" {
		logger.Info("Initial search term: " + st)
	} else {
		logger.Info("No initial search term provided.")
	}
	return st
}
