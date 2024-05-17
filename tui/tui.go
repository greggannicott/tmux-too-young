package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	message     string
	keys        keyMap
	help        help.Model
	projects    []projectModel
	cursorIndex int
}

type keyMap struct {
	Quit key.Binding
	Up   key.Binding
	Down key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Quit},
		{},
	}
}

var DefaultKeyMap = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Up: key.NewBinding(
		key.WithKeys("ctrl+k", "up", "ctrl+p"),
		key.WithHelp("ctrl+k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("ctrl+j", "down", "ctrl+n"),
		key.WithHelp("ctrl+j", "down"),
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
		case key.Matches(msg, DefaultKeyMap.Up):
			if m.cursorIndex == 0 {
				m.cursorIndex = len(m.projects) - 1
			} else {
				m.cursorIndex--
			}
		case key.Matches(msg, DefaultKeyMap.Down):
			if m.cursorIndex == len(m.projects)-1 {
				m.cursorIndex = 0
			} else {
				m.cursorIndex++
			}
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
	for i, p := range m.projects {
		cursor := " "
		if i == m.cursorIndex {
			cursor = ">"
		}
		sb.WriteString(fmt.Sprintf("%v %v\n", cursor, p.fullPath))
	}
	sb.WriteString("\n")
	sb.WriteString(m.help.View(m.keys))
	return sb.String()
}
