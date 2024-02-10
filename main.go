package main

import (
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

var Version = "Development"

type worktreeDetails struct {
	worktree string
	branch   string
}

var projects []project

// To run in terminal: go run tmux-too-young
// Config file: ~/.tmux-too-young.yaml
func main() {
	app := &cli.App{
		Name:            "tmux-too-young",
		Usage:           "The Very Special tmux Session Opener...",
		HideHelpCommand: true,
		Version:         Version,
		Action: func(*cli.Context) error {
			ensureAppCanRun()
			initialSearchTerm := getInitialSearchTerm()
			config := getConfig()
			scanProjectDirectories(config)
			selectedProject := getSelectionFromUser(initialSearchTerm)
			launchProject(selectedProject)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getInitialSearchTerm() string {
	return strings.Join(os.Args[1:], " ")
}