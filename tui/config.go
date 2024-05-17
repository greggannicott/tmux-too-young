package tui

import (
	"tmux-too-young/config"

	tea "github.com/charmbracelet/bubbletea"
)

type ConfigRetrievedMsg struct {
	searchDirectories []string
}

func InitConfig() tea.Msg {
	c := config.GetConfig()

	return ConfigRetrievedMsg{
		searchDirectories: c.SearchDirectories,
	}
}
