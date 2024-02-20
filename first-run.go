package main

import (
	"fmt"
)

func ensureAppCanRun() {
	if configNeedsCreating() {
		sd := getSearchDirectoriesFromUser()
		createConfig(sd)
		displayConfigInstructions()
	}
}

func configNeedsCreating() bool {
	return !configExists()
}

func getSearchDirectoriesFromUser() string {
	var sd string
	fmt.Println("Please provide a comma separated list of directories containing projects you wish to open with tmux-too-young.")
	fmt.Println()
	fmt.Println("For example, if you have a collection of projects inside ~/code/, and folders containing repos inside ~/, enter \"~/code/,~/\":")
	fmt.Println()
	fmt.Print("> ")
	fmt.Scan(&sd)
	return sd
}

func displayConfigInstructions() {
	fmt.Println()
	fmt.Println("Thanks!")
	fmt.Println()
	fmt.Println("`.tmux-too-young.yaml` has been created in your home directory.")
	fmt.Println("")
	fmt.Println("You can update this file when you wish to add/remove search directories.")
	fmt.Println("")
	fmt.Println("PRESS ANY KEY")
	fmt.Scanln()
}
