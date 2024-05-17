package tui

import (
	proj "tmux-too-young/project"

	tea "github.com/charmbracelet/bubbletea"
)

type projectModel struct {
	fullPath string
}

type ProjectsRetrievedMsg struct {
	projects []projectModel
}

func getProjects(sd []string) tea.Cmd {
	return func() tea.Msg {
		projects := proj.ScanProjectDirectories(sd)
		var projectModels []projectModel
		for _, v := range projects {
			projectModels = append(projectModels, projectModel{fullPath: v.FullPath})
		}
		return ProjectsRetrievedMsg{
			projects: projectModels,
		}
	}
}
