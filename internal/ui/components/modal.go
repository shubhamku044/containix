package components

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors used across the application
	primaryColor   = lipgloss.Color("39")  // Light blue
	secondaryColor = lipgloss.Color("62")  // Purple
	accentColor    = lipgloss.Color("226") // Yellow
	subtleColor    = lipgloss.Color("240") // Gray

	// Base styles
	modalTitleStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Padding(0, 1)

	modalBorderStyle = lipgloss.RoundedBorder()

	// Help style for showing keyboard shortcuts
	helpStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			Italic(true)
)

type ModalModel struct {
	viewport    viewport.Model
	content     string
	width       int
	height      int
	parentModel tea.Model
	title       string
}

func NewModal(title string, content string, width, height int, parentModel tea.Model) ModalModel {
	// Adjust viewport to account for borders, padding, and title bar
	vp := viewport.New(width-6, height-10)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(subtleColor)
	vp.SetContent(content)

	return ModalModel{
		viewport:    vp,
		content:     content,
		width:       width,
		height:      height,
		parentModel: parentModel,
		title:       title,
	}
}

// Init implements tea.Model
func (m ModalModel) Init() tea.Cmd {
	// No initialization needed, just return nil
	return nil
}

func (m ModalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m.parentModel, nil
		case "j", "down":
			m.viewport.LineDown(1)
		case "k", "up":
			m.viewport.LineUp(1)
		case "g", "home":
			m.viewport.GotoTop()
		case "G", "end":
			m.viewport.GotoBottom()
		case "f", "pagedown":
			m.viewport.HalfViewDown()
		case "b", "pageup":
			m.viewport.HalfViewUp()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 6
		m.viewport.Height = msg.Height - 10
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m ModalModel) View() string {
	// Determine modal dimensions - make it feel more spacious
	modalWidth := m.width * 85 / 100
	modalHeight := m.height * 80 / 100

	// Create a header with the title
	header := modalTitleStyle.Render(m.title)

	// Add a scrollbar indicator and info text to the viewport
	viewportContent := m.viewport.View()

	// Help text at the bottom
	helpText := helpStyle.Render("↑/k: up • ↓/j: down • g/home: top • G/end: bottom • q/esc: close")

	// Build the complete modal content
	modalContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		viewportContent,
		"",
		helpText,
	)

	// Style for the entire modal
	modalStyle := lipgloss.NewStyle().
		Width(modalWidth).
		Height(modalHeight).
		BorderStyle(modalBorderStyle).
		BorderForeground(secondaryColor).
		Padding(1, 2)

	// Render the modal content inside the styled box
	content := modalStyle.Render(modalContent)

	// Center the modal in available space
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
		lipgloss.WithWhitespaceChars(""),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("235")),
	)
}
