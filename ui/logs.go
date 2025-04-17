package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type logViewModel struct {
	viewport viewport.Model
}

func newLogView() logViewModel {
	vp := viewport.New(50, 20)
	vp.SetContent("← Select a container to view logs here.")
	return logViewModel{viewport: vp}
}

func (m logViewModel) Update(msg tea.Msg) (logViewModel, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m logViewModel) View() string {
	return lipgloss.NewStyle().Padding(1).Width(50).Render(m.viewport.View())
}
