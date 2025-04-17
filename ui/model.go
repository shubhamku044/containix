package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	containers containerListModel
	logs       logViewModel
	focusLeft  bool
}

func NewModel() tea.Model {
	return model{
		containers: newContainerList(),
		logs:       newLogView(),
		focusLeft:  true,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.focusLeft = !m.focusLeft
		case "q":
			return m, tea.Quit
		}

		if m.focusLeft {
			m.containers, cmd = m.containers.Update(msg)
		} else {
			m.logs, cmd = m.logs.Update(msg)
		}
	}

	return m, cmd
}

func (m model) View() string {
	left := m.containers.View()
	right := m.logs.View()

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
