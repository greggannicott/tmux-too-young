package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type projectDirectory struct {
	fullPath string
}

func main() {
	projectDirectory := getDirectories()
	selectedProject := getSelectionFromFzf(projectDirectory)
	fmt.Println("Selection", selectedProject)
}

func getDirectories() []projectDirectory {
	launchableDirs := []projectDirectory{}
	// The following is hard coded but eventually will be obtained via a loop over a config entry.
	currentDir := "/Users/greggannicott/code/"
	dirs, _ := os.ReadDir(currentDir)
	for _, dir := range dirs {
		fullPath := currentDir + dir.Name()
		_, err := os.Stat(fullPath + "/.git/")
		if err == nil {
			launchableDir := projectDirectory{
				fullPath: fullPath,
			}
			launchableDirs = append(launchableDirs, launchableDir)
		}
	}
	return launchableDirs
}

func getSelectionFromFzf(choices []projectDirectory) string {
	var input string
	for _, choice := range choices {
		input += choice.fullPath + "\n"
	}
	cmd := exec.Command("fzf-tmux", "-p", "--cycle", "--reverse", "--border", "--info=inline-right", "--header=Select a Project to open in tmux:")
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(stdout.String())
}

func (l projectDirectory) getFriendlyName() string {
	return l.fullPath
}
