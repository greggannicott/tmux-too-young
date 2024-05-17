package project

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func ScanProjectDirectories(sd []string) []Project {
	for i := 0; i < len(sd); i++ {
		scanProjectDirectoriesFor(sd[i])
	}
	return projects
}

func scanProjectDirectoriesFor(rootDir string) {
	rootDir = prepareRootDir(rootDir)
	dirs, _ := os.ReadDir(rootDir)
	for _, dir := range dirs {
		basePath := rootDir + dir.Name()
		_, gitDirErr := os.Stat(basePath + "/.git/")
		_, gitFileErr := os.Stat(basePath + "/.git")
		_, tmuxpFileErr := os.Stat(basePath + "/.tmuxp.yaml")
		_, tmuxTooYoungFileErr := os.Stat(basePath + "/.tmux-too-young")
		if gitDirErr == nil || gitFileErr == nil || tmuxpFileErr == nil || tmuxTooYoungFileErr == nil {
			// Find out if there are any worktrees in this directory
			worktrees := getWorktreesForProject(basePath)
			if projectHasWorktrees(worktrees, basePath) {
				for _, w := range worktrees {
					projectHasTmuxpFile := projectHasTmuxpFile(basePath + "/" + w.branch)
					launchableDir := Project{
						basePath:      basePath,
						FullPath:      basePath + "/" + w.branch,
						isWorktree:    true,
						branch:        w.branch,
						supportsTmuxp: projectHasTmuxpFile,
					}
					projects = append(projects, launchableDir)
				}
			} else {
				projectHasTmuxpFile := projectHasTmuxpFile(basePath)
				launchableDir := Project{
					basePath:      basePath,
					FullPath:      basePath,
					isWorktree:    false,
					supportsTmuxp: projectHasTmuxpFile,
				}
				projects = append(projects, launchableDir)
			}
		}
	}
}

func prepareRootDir(rd string) string {
	rd = addTrailigSlash(rd)
	rd = replaceTildeWithHomeDirectory(rd)
	return rd
}

func addTrailigSlash(s string) string {
	if !strings.HasSuffix(s, string(os.PathSeparator)) {
		s = s + string(os.PathSeparator)
	}
	return s
}

func replaceTildeWithHomeDirectory(s string) string {
	userHomeDir, _ := os.UserHomeDir()
	return strings.Replace(s, "~", userHomeDir, 1)
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
			// Capture everything following the first space (including other spaces)
			value := strings.Join(keyValue[1:], " ")
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
