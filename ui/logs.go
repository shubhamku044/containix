package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type logViewModel struct {
	viewport viewport.Model
	width    int
	height   int
}

func newLogView() logViewModel {
	vp := viewport.New(50, 20)
	vp.SetContent("← Select a container to view logs here.")

	return logViewModel{
		viewport: vp,
		width:    50,
		height:   20,
	}
}

func (m *logViewModel) SetContent(content string) {
	m.viewport.SetContent(content)
	// Reset viewport to top when content changes
	m.viewport.GotoTop()
}

func (m logViewModel) Update(msg tea.Msg) (logViewModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width / 2
		m.height = msg.Height - 2
		m.viewport.Width = m.width - 4
		m.viewport.Height = m.height - 2
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m logViewModel) View() string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1).
		Width(m.width).
		Height(m.height).
		Render(m.viewport.View())
}
