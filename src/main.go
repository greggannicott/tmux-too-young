package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
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
	populateProjectDirectories(config)
	selectedProject := getSelectionFromFzf(initialSearchTerm)
	openProject(selectedProject)
}

func getInitialSearchTerm() string {
	return strings.Join(os.Args[1:], " ")
}

func populateProjectDirectories(config config) {
	for i := 0; i < len(config.SearchDirectories); i++ {
		populateProjectDirectoriesFor(config.SearchDirectories[i])
	}
}

func populateProjectDirectoriesFor(rootDir string) {
	dirs, _ := os.ReadDir(rootDir)
	for _, dir := range dirs {
		basePath := rootDir + dir.Name()
		_, dirErr := os.Stat(basePath + "/.git/")
		_, fileErr := os.Stat(basePath + "/.git")
		if dirErr == nil || fileErr == nil {
			// Find out if there are any worktrees in this directory
			worktrees := getWorktrees(basePath)
			if projectHasWorktrees(worktrees, basePath) {
				for _, w := range worktrees {
					projectHasTmuxpFile := projectHasTmuxpFile(basePath + "/" + w.branch)
					launchableDir := project{
						basePath:      basePath,
						fullPath:      basePath + "/" + w.branch,
						isWorktree:    true,
						branch:        w.branch,
						supportsTmuxp: projectHasTmuxpFile,
					}
					projects = append(projects, launchableDir)
				}
			} else {
				projectHasTmuxpFile := projectHasTmuxpFile(basePath)
				launchableDir := project{
					basePath:      basePath,
					fullPath:      basePath,
					isWorktree:    false,
					supportsTmuxp: projectHasTmuxpFile,
				}
				projects = append(projects, launchableDir)
			}
		}
	}
}

func getWorktrees(basePath string) []worktreeDetails {
	var worktrees []worktreeDetails

	cmd := exec.Command("git", "-C", basePath, "worktree", "list", "--porcelain")
	commandOutput, _ := cmd.Output()

	// Break the results down into "rawWorktrees".
	rawWorktrees := strings.Split(string(commandOutput), "\n\n")

	// Iterate over each rawWorktree and parse them
	for _, rawWorktree := range rawWorktrees {
		var worktree worktreeDetails

		for _, line := range strings.Split(string(rawWorktree), "\n") {
			keyValue := strings.Split(line, " ")
			// If we have a key but no value, we're not interested..
			if len(keyValue) == 1 {
				continue
			}
			key := keyValue[0]
			value := keyValue[1]
			if key == "worktree" {
				worktree.worktree = value
			} else if key == "branch" {
				// Strip away the cruft and just leave the branch name
				re := regexp.MustCompile(`refs\/heads\/(.*)`)
				matches := re.FindStringSubmatch(value)
				if len(matches) > 1 {
					worktree.branch = string(matches[1])
				} else {
					worktree.branch = value
				}
				// If this is the branch key, then its the final key for the worktree.
				// As a result we need to add the worktree to the slice
				worktrees = append(worktrees, worktree)
			}
		}
	}

	return worktrees
}

func projectHasWorktrees(worktreeDetails []worktreeDetails, basePath string) bool {
	if len(worktreeDetails) == 0 {
		return false
	}
	return worktreeDetails[0].worktree != basePath
}

func projectHasTmuxpFile(basePath string) bool {
	tmuxpPath := basePath + "/.tmuxp.yaml"
	_, err := os.Stat(tmuxpPath)
	return err == nil
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
