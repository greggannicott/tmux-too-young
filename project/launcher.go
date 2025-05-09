package project

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func LaunchProject(selectedProject project) {
	if selectedProject.supportsTmuxp == true {
		launchProjectUsingTmuxp(selectedProject)
	} else {
		if isInsideOfTmux() {
			if sessionIsUnderway(selectedProject) {
				attachToProjectFromWithinTmux(selectedProject)
			} else {
				launchProjectFromWithinTmux(selectedProject)
			}
		} else {
			if sessionIsUnderway(selectedProject) {
				attachToProjectFromOutsideOfTmux(selectedProject)
			} else {
				launchProjectFromOutsideOfTmux(selectedProject)
			}
		}
	}
}

func launchProjectFromWithinTmux(p project) {
	newSessionCmd := exec.Command("tmux", "new-session", "-d", "-s", p.getSessionName(), "-c", p.fullPath)
	err := newSessionCmd.Start()
	if err != nil {
		fmt.Println("Error creating new tmux session:", err)
	}
	waitForSessionToBeReady(p)
	attachToProjectFromWithinTmux(p)
}

func waitForSessionToBeReady(p project) {
	maxTries := 5
	for i := 0; i < maxTries; i++ {
		if sessionIsUnderway(p) {
			break
		}
		fmt.Printf("Attempted to open new session [%v] but it isn't ready yet...", p.getFriendlyName())
		time.Sleep(500 * time.Millisecond)
		if i == maxTries-1 {
			fmt.Println("Giving up on opening new session.")
			return
		}
	}
}

func launchProjectFromOutsideOfTmux(p project) {
	cmd := exec.Command("tmux", "new-session", "-s", p.getSessionName(), "-c", p.fullPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error creation new session:", err)
	}
}

func launchProjectUsingTmuxp(p project) {
	cmd := exec.Command("tmuxp", "load", p.getTmuxpPath(), "-s", p.getSessionName(), "-y")
	cmd.Stdin = os.Stdin
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
