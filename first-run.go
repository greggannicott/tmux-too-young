package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ensureAppCanRun() {
	if configNeedsCreating() {
		sd := getSearchDirectoriesFromUser()
		createConfig(sd)
		displayConfigInstructions(sd)
	}
}

func configNeedsCreating() bool {
	return !configExists()
}

func getSearchDirectoriesFromUser() []string {
	var searchDirectory string
	var searchDirectories []string
	for len(searchDirectories) == 0 || searchDirectory != "" {
		displayInputInstructions(searchDirectories)
		searchDirectory = ""
		fmt.Scanln(&searchDirectory)
		searchDirectory = strings.TrimSpace(searchDirectory)
		if searchDirectory != "" {
			searchDirectories = append(searchDirectories, searchDirectory)
		}
	}
	return searchDirectories
}

func displayInputInstructions(searchDirectories []string) {
	clearScreen()
	fmt.Println("Please enter a directory you would like tmux-too-young to scan for projects.")
	fmt.Println()
	fmt.Println("For example, if you have a collection of projects inside ~/code/, enter \"~/code/\".")
	fmt.Println()
	if len(searchDirectories) == 0 {
		fmt.Println("At least one entry is required.")
	} else {
		fmt.Println("Press enter without entering a value to continue.")
	}
	fmt.Println()
	if len(searchDirectories) == 0 {
		fmt.Println("Note: You will be prompted for additional directories one you've entered this one.")
	} else {
		fmt.Printf("Existing Search Directories: [%s]\n", strings.Join(searchDirectories, ", "))
	}
	fmt.Println()
	fmt.Print("> ")
}

func displayConfigInstructions(searchDirectories []string) {
	clearScreen()
	fmt.Println("Thanks!")
	fmt.Println()
	fmt.Println("`.tmux-too-young.yaml` has been created in your home directory.")
	fmt.Println()
	fmt.Println("The following search directories will be used:")
	fmt.Println("")
	for _, directory := range searchDirectories {
		fmt.Printf("* %s\n", directory)
	}
	fmt.Println()
	fmt.Println("You can update this file when you wish to add/remove search directories.")
	fmt.Println("")
	fmt.Println("PRESS ANY KEY")
	fmt.Scanln()
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
