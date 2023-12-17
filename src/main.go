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

var projectDirectories []projectDirectory

func main() {
	getProjectDirectories()
	projectDirectory := getSelectionFromFzf()
	if isInsideOfTmux() {
		if sessionIsUnderway(projectDirectory) {
			attachToProjectFromWithinTmux(projectDirectory)
		} else {
			openProjectFromWithinTmux(projectDirectory)
		}
	} else {
		fmt.Println("Unable to open session as we are not inside of Tmux...")
	}
}

func getProjectDirectories() {
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
			projectDirectories = append(projectDirectories, launchableDir)
		}
	}
}

func getSelectionFromFzf() projectDirectory {
	var input string
	for _, choice := range projectDirectories {
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
	selectedProjectName := strings.TrimSpace(stdout.String())
	return findProjectDirectoryByFriendlyName(selectedProjectName)
}

func openProjectFromWithinTmux(projectDirectory projectDirectory) {
	newSessionCmd := exec.Command("tmux", "new-session", "-d", "-s", projectDirectory.getSessionName(), "-c", "/users/greggannicott/code/tmux-too-young")
	err := newSessionCmd.Start()
	if err != nil {
		fmt.Println("Error creating new tmux session:", err)
	}
	attachToProjectFromWithinTmux(projectDirectory)
}

func attachToProjectFromWithinTmux(projectDirectory projectDirectory) {
	switchSessionCmd := exec.Command("tmux", "switch-client", "-t", projectDirectory.getSessionName())
	switchSessionCmdErr := switchSessionCmd.Start()
	if switchSessionCmdErr != nil {
		fmt.Println("Error switching to new session:", switchSessionCmdErr)
	}
}

func findProjectDirectoryByFriendlyName(name string) projectDirectory {
	var matchingProjectDirectory projectDirectory
	for _, projectDirectory := range projectDirectories {
		if projectDirectory.getFriendlyName() == name {
			matchingProjectDirectory = projectDirectory
			break
		}
	}
	return matchingProjectDirectory
}

func sessionIsUnderway(projectDirectory projectDirectory) bool {
	sessionIsUserwayCmd := exec.Command("/bin/zsh", "-c", "tmux list-sessions | cut -d ':' -f 1 | grep "+projectDirectory.getSessionName())
	output, _ := sessionIsUserwayCmd.Output()
	return string(output) != ""
}

func isInsideOfTmux() bool {
	_, isInsideOfTmux := os.LookupEnv("TMUX")
	return isInsideOfTmux
}

func (l projectDirectory) getFriendlyName() string {
	return l.fullPath
}

func (p projectDirectory) getSessionName() string {
	fileInfo, _ := os.Stat(p.fullPath)
	return fileInfo.Name()
}
