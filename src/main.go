package main

import (
	"os"
	"strings"
)

type worktreeDetails struct {
	worktree string
	branch   string
}

var projects []project

// To run in terminal: go run tmux-too-young
// Config file: ~/.tmux-too-young.yaml
func main() {
	initialSearchTerm := getInitialSearchTerm()
	config := getConfig()
	scanProjectDirectories(config)
	selectedProject := getSelectionFromUser(initialSearchTerm)
	launchProject(selectedProject)
}

func getInitialSearchTerm() string {
	return strings.Join(os.Args[1:], " ")
}
