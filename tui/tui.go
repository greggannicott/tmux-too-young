package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	message  string
	keys     keyMap
	help     help.Model
	projects []projectModel
}

type keyMap struct {
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit},
		{},
	}
}

var DefaultKeyMap = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func Display() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error displaying TUI: %v\n\n", err)
	}
}
func initModel() tea.Model {
	return model{
		message: "",
		help:    help.New(),
		keys:    DefaultKeyMap,
	}
}

func (m model) Init() tea.Cmd {
	return InitConfig
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit
		}
	case ConfigRetrievedMsg:
		return m, getProjects(msg.searchDirectories)
	case ProjectsRetrievedMsg:
		m.projects = msg.projects
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("You've done too much... tmux-too-young...\n")
	sb.WriteString("\n")
	for _, p := range m.projects {
		sb.WriteString(fmt.Sprintf("* %v\n", p.fullPath))
	}
	sb.WriteString("\n")
	sb.WriteString(m.help.View(m.keys))
	return sb.String()
}
