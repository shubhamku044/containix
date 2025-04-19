package components

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ModalModel represents a modal dialog
type ModalModel struct {
	viewport    viewport.Model
	content     string
	width       int
	height      int
	parentModel tea.Model
}

// NewModal creates a new modal
func NewModal(content string, width, height int, parentModel tea.Model) ModalModel {
	vp := viewport.New(width-4, height-6)
	vp.SetContent(content)

	return ModalModel{
		viewport:    vp,
		content:     content,
		width:       width,
		height:      height,
		parentModel: parentModel,
	}
}

// Update updates the model
func (m ModalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m.parentModel, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 6
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the model
func (m ModalModel) View() string {
	modalWidth := m.width * 80 / 100
	modalHeight := m.height * 70 / 100

	modalStyle := lipgloss.NewStyle().
		Width(modalWidth).
		Height(modalHeight).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2)

	content := modalStyle.Render(m.viewport.View())
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

