package views

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LogViewModel struct {
	viewport viewport.Model
}

func NewLogViewModel() LogViewModel {
	vp := viewport.New(50, 20)
	vp.SetContent("← Select a container to view logs here.")
	return LogViewModel{viewport: vp}
}

func (m *LogViewModel) SetContent(content string) {
	m.viewport.SetContent(content)
}

func (m LogViewModel) Update(msg tea.Msg) (LogViewModel, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m LogViewModel) View() string {
	return lipgloss.NewStyle().Padding(1).Width(50).Render(m.viewport.View())
}

