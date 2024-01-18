package main

import (
	"fmt"
	"os"
	"os/exec"
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
	selectedProject := getSelectionFromFzf(initialSearchTerm)
	openProject(selectedProject)
}

func getInitialSearchTerm() string {
	return strings.Join(os.Args[1:], " ")
}

func openProject(selectedProject project) {
	if selectedProject.supportsTmuxp == true {
		openProjectUsingTmuxp(selectedProject)
	} else {
		if isInsideOfTmux() {
			if sessionIsUnderway(selectedProject) {
				attachToProjectFromWithinTmux(selectedProject)
			} else {
				openProjectFromWithinTmux(selectedProject)
			}
		} else {
			if sessionIsUnderway(selectedProject) {
				attachToProjectFromOutsideOfTmux(selectedProject)
			} else {
				openProjectFromOutsideOfTmux(selectedProject)
			}
		}
	}
}

func openProjectFromWithinTmux(p project) {
	newSessionCmd := exec.Command("tmux", "new-session", "-d", "-s", p.getSessionName(), "-c", p.fullPath)
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

func openProjectUsingTmuxp(p project) {
	cmd := exec.Command("tmuxp", "load", p.getTmuxpPath(), "-s", p.getSessionName(), "-y")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running tmuxp:", err)
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

func sessionIsUnderway(p project) bool {
	sessionIsUserwayCmd := exec.Command("/bin/zsh", "-c", "tmux list-sessions | cut -d ':' -f 1 | grep '"+p.getSessionName()+"'")
	output, _ := sessionIsUserwayCmd.Output()
	return string(output) != ""
}

func isInsideOfTmux() bool {
	_, isInsideOfTmux := os.LookupEnv("TMUX")
	return isInsideOfTmux
}
