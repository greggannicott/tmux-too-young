package main

import (
	"log"
	"os"
	config "tmux-too-young/config"
	first_run "tmux-too-young/first-run"
	project "tmux-too-young/project"

	"github.com/urfave/cli/v2"
)

var Version = "Development"

// To run in terminal: go run tmux-too-young
// Config file: ~/.tmux-too-young.yaml
func main() {
	var initialSearchTerm string
	app := &cli.App{
		Name:            "tmux-too-young",
		Usage:           "The Very Special tmux Session Opener...",
		HideHelpCommand: true,
		Version:         Version,
		Commands: []*cli.Command{
			{
				Name:  "open",
				Usage: "Open a tmux session",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "search",
						Usage:       "Initial search term.",
						Destination: &initialSearchTerm,
						Aliases:     []string{"s"},
					},
				},
				Action: func(cCtx *cli.Context) error {
					first_run.EnsureAppCanRun()
					c := config.GetConfig()
					project.ScanProjectDirectories(c)
					selectedProject := project.GetSelectionFromUser(initialSearchTerm)
					project.LaunchProject(selectedProject)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
