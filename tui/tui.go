package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	message string
}

func Display() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error displaying TUI: %v\n\n", err)
	}
}
func initModel() tea.Model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	output := "\n"
	output += "tmux-too-young...\n"
	output += "\n"
	output += "q to quit\n"
	if len(m.message) > 0 {
		output += m.message
	}
	return output
}
