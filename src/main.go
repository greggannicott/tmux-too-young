package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type launchableDir struct {
	fullPath string
}

func main() {
	launchableDirs := getDirs()
	selection := getSelectionFromFzf(launchableDirs)
	fmt.Println("Selection", selection)
}

func getDirs() []launchableDir {
	launchableDirs := []launchableDir{}
	// The following is hard coded but eventually will be obtained via a loop over a config entry.
	currentDir := "/Users/greggannicott/code/"
	dirs, _ := os.ReadDir(currentDir)
	for _, dir := range dirs {
		fullPath := currentDir + dir.Name()
		_, err := os.Stat(fullPath + "/.git/")
		if err == nil {
			launchableDir := launchableDir{
				fullPath: fullPath,
			}
			launchableDirs = append(launchableDirs, launchableDir)
		}
	}
	return launchableDirs
}

func getSelectionFromFzf(choices []launchableDir) string {
	var input string
	for _, choice := range choices {
		input += choice.fullPath + "\n"
	}
	cmd := exec.Command("fzf-tmux", "-p", "--cycle", "--reverse", "--border", "--info=inline-right", "--header=Select a repo to open in tmux:")
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

func (l launchableDir) getFriendlyName() string {
	return l.fullPath
}
