package main

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func scanProjectDirectories(config config) {
	for i := 0; i < len(config.SearchDirectories); i++ {
		scanProjectDirectoriesFor(config.SearchDirectories[i])
	}
}

func scanProjectDirectoriesFor(rootDir string) {
	dirs, _ := os.ReadDir(rootDir)
	for _, dir := range dirs {
		basePath := rootDir + dir.Name()
		_, dirErr := os.Stat(basePath + "/.git/")
		_, fileErr := os.Stat(basePath + "/.git")
		if dirErr == nil || fileErr == nil {
			// Find out if there are any worktrees in this directory
			worktrees := getWorktreesForProject(basePath)
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

func getWorktreesForProject(basePath string) []worktreeDetails {
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
