package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type project struct {
	fullPath string
}

var projects []project

func main() {
	getProjectDirectories()
	project := getSelectionFromFzf()
	if isInsideOfTmux() {
		if sessionIsUnderway(project) {
			attachToProjectFromWithinTmux(project)
		} else {
			openProjectFromWithinTmux(project)
		}
	} else {
		if sessionIsUnderway(project) {
			attachToProjectFromOutsideOfTmux(project)
		} else {
			openProjectFromOutsideOfTmux(project)
		}
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
			launchableDir := project{
				fullPath: fullPath,
			}
			projects = append(projects, launchableDir)
		}
	}
}

func getSelectionFromFzf() project {
	var input string
	for _, choice := range projects {
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

func openProjectFromWithinTmux(p project) {
	newSessionCmd := exec.Command("tmux", "new-session", "-d", "-s", p.getSessionName(), "-c", "/users/greggannicott/code/tmux-too-young")
	err := newSessionCmd.Start()
	if err != nil {
		fmt.Println("Error creating new tmux session:", err)
	}
	attachToProjectFromWithinTmux(p)
}

func openProjectFromOutsideOfTmux(p project) {
	cmd := exec.Command("tmux", "new-session", "-s", p.getSessionName(), "-c", p.fullPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error creation new session:", err)
	}
}

func attachToProjectFromWithinTmux(p project) {
	cmd := exec.Command("tmux", "switch-client", "-t", p.getSessionName())
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error switching to session:", err)
	}
}

func attachToProjectFromOutsideOfTmux(p project) {
	cmd := exec.Command("tmux", "attach-session", "-t", p.getSessionName())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error attaching to exiting session:", err)
	}
}

func findProjectDirectoryByFriendlyName(name string) project {
	var matchingProjectDirectory project
	for _, p := range projects {
		if p.getFriendlyName() == name {
			matchingProjectDirectory = p
			break
		}
	}
	return matchingProjectDirectory
}

func sessionIsUnderway(p project) bool {
	sessionIsUserwayCmd := exec.Command("/bin/zsh", "-c", "tmux list-sessions | cut -d ':' -f 1 | grep "+p.getSessionName())
	output, _ := sessionIsUserwayCmd.Output()
	return string(output) != ""
}

func isInsideOfTmux() bool {
	_, isInsideOfTmux := os.LookupEnv("TMUX")
	return isInsideOfTmux
}

func (l project) getFriendlyName() string {
	return l.fullPath
}

func (p project) getSessionName() string {
	fileInfo, _ := os.Stat(p.fullPath)
	return fileInfo.Name()
}
