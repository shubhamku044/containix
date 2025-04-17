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
	containers, err := newContainerList()
	if err != nil {
		panic(err)
	}
	return model{
		containers: containers,
		logs:       newLogView(),
		focusLeft:  true,
	}
}

func (m model) Init() tea.Cmd {
	return m.containers.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.focusLeft = !m.focusLeft
			return m, nil
		case "q":
			return m, tea.Quit
		}
	}

	if m.focusLeft {
		var newModel tea.Model
		newModel, cmd = m.containers.Update(msg)
		m.containers = newModel.(containerListModel)
	}

	return m, cmd
}

func (m model) View() string {
	left := m.containers.View()
	right := m.logs.View()
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
