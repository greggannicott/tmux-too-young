package project

import (
	"os"
	"strings"
)

type Project struct {
	basePath      string
	FullPath      string
	isWorktree    bool
	branch        string
	supportsTmuxp bool
}

type worktreeDetails struct {
	worktree string
	branch   string
}

var projects []Project

func (l Project) getFriendlyName() string {
	if l.isWorktree {
		return l.basePath + " -> " + l.branch
	} else {
		return l.FullPath
	}
}

func (p Project) getSessionName() string {
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
func (p Project) getTmuxpPath() string {
	if p.isWorktree {
		return p.FullPath + "/.tmuxp.yaml"
	} else {
		return p.basePath + "/.tmuxp.yaml"
	}
}
