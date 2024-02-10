package main

import (
	"os"
	"strings"
)

type project struct {
	basePath      string
	fullPath      string
	isWorktree    bool
	branch        string
	supportsTmuxp bool
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
func (p project) getTmuxpPath() string {
	if p.isWorktree {
		return p.fullPath + "/.tmuxp.yaml"
	} else {
		return p.basePath + "/.tmuxp.yaml"
	}
}
