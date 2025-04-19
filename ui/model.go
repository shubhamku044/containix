package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	containers containerListModel
	logs       logViewModel
	focusLeft  bool
	width      int
	height     int
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update window size for both components
		m.width = msg.Width
		m.height = msg.Height

		containerMsg := tea.WindowSizeMsg{
			Width:  msg.Width / 2,
			Height: msg.Height,
		}

		logMsg := tea.WindowSizeMsg{
			Width:  msg.Width / 2,
			Height: msg.Height,
		}

		var containerModel tea.Model
		containerModel, cmd = m.containers.Update(containerMsg)
		m.containers = containerModel.(containerListModel)
		cmds = append(cmds, cmd)

		m.logs, cmd = m.logs.Update(logMsg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case LogMessage:
		// When logs are fetched, update the log view
		m.logs.SetContent(msg.Logs)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.focusLeft = !m.focusLeft
			return m, nil
		case "q":
			if !m.focusLeft {
				// If focus is on logs view, just return focus to container list
				m.focusLeft = true
				return m, nil
			}
			return m, tea.Quit
		}
	}

	// Route messages to the appropriate component based on focus
	if m.focusLeft {
		var containerModel tea.Model
		containerModel, cmd = m.containers.Update(msg)
		m.containers = containerModel.(containerListModel)
	} else {
		m.logs, cmd = m.logs.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	left := m.containers.View()
	right := m.logs.View()

	// Apply styles based on focus
	leftStyle := lipgloss.NewStyle()
	rightStyle := lipgloss.NewStyle()

	if m.focusLeft {
		leftStyle = leftStyle.BorderForeground(lipgloss.Color("62"))
	} else {
		rightStyle = rightStyle.BorderForeground(lipgloss.Color("62"))
	}

	styledLeft := leftStyle.Render(left)
	styledRight := rightStyle.Render(right)

	return lipgloss.JoinHorizontal(lipgloss.Top, styledLeft, styledRight)
}
