package project

import (
	"fmt"
	"os"
	"os/exec"
)

func LaunchProject(selectedProject Project) {
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

func launchProjectFromWithinTmux(p Project) {
	newSessionCmd := exec.Command("tmux", "new-session", "-d", "-s", p.getSessionName(), "-c", p.FullPath)
	err := newSessionCmd.Start()
	if err != nil {
		fmt.Println("Error creating new tmux session:", err)
	}
	attachToProjectFromWithinTmux(p)
}

func launchProjectFromOutsideOfTmux(p Project) {
	cmd := exec.Command("tmux", "new-session", "-s", p.getSessionName(), "-c", p.FullPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error creation new session:", err)
	}
}

func launchProjectUsingTmuxp(p Project) {
	cmd := exec.Command("tmuxp", "load", p.getTmuxpPath(), "-s", p.getSessionName(), "-y")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running tmuxp:", err)
	}
}

func attachToProjectFromWithinTmux(p Project) {
	cmd := exec.Command("tmux", "switch-client", "-t", p.getSessionName())
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error switching to session:", err)
	}
}

func attachToProjectFromOutsideOfTmux(p Project) {
	cmd := exec.Command("tmux", "attach-session", "-t", p.getSessionName())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error attaching to exiting session:", err)
	}
}

func sessionIsUnderway(p Project) bool {
	sessionIsUserwayCmd := exec.Command("/bin/zsh", "-c", "tmux list-sessions | cut -d ':' -f 1 | grep '"+p.getSessionName()+"'")
	output, _ := sessionIsUserwayCmd.Output()
	return string(output) != ""
}

func isInsideOfTmux() bool {
	_, isInsideOfTmux := os.LookupEnv("TMUX")
	return isInsideOfTmux
}
