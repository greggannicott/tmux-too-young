package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type project struct {
	basePath   string
	fullPath   string
	isWorktree bool
	branch     string
}
type worktreeDetails struct {
	worktree string
	branch   string
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
		basePath := currentDir + dir.Name()
		_, dirErr := os.Stat(basePath + "/.git/")
		_, fileErr := os.Stat(basePath + "/.git")
		if dirErr == nil || fileErr == nil {
			// Find out if there are any worktrees in this directory
			worktrees := getWorktrees(basePath)
			if projectHasWorktrees(worktrees, basePath) {
				for _, w := range worktrees {
					launchableDir := project{
						basePath:   basePath,
						fullPath:   basePath + "/" + w.branch,
						isWorktree: true,
						branch:     w.branch,
					}
					projects = append(projects, launchableDir)
				}
			} else {
				launchableDir := project{
					basePath:   basePath,
					fullPath:   basePath,
					isWorktree: false,
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
	return worktreeDetails[0].worktree != basePath
}

func getSelectionFromFzf() project {
	var input string
	for _, choice := range projects {
		input += choice.getFriendlyName() + "\n"
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
	sessionIsUserwayCmd := exec.Command("/bin/zsh", "-c", "tmux list-sessions | cut -d ':' -f 1 | grep '"+p.getSessionName()+"'")
	output, _ := sessionIsUserwayCmd.Output()
	return string(output) != ""
}

func isInsideOfTmux() bool {
	_, isInsideOfTmux := os.LookupEnv("TMUX")
	return isInsideOfTmux
}

func (l project) getFriendlyName() string {
	if l.isWorktree {
		return l.basePath + " -> " + l.branch
	} else {
		return l.fullPath
	}
}

func (p project) getSessionName() string {
	fileInfo, _ := os.Stat(p.basePath)
	// `.`s need to be replaced as they're not allowed in a tmux name
	safeName := strings.ReplaceAll(fileInfo.Name(), ".", "_")
	if p.isWorktree {
		safeBranch := strings.ReplaceAll(p.branch, ".", "_")
		return safeName + " -> " + safeBranch
	} else {
		return safeName
	}
}
