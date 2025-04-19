package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shubhamku044/containix/internal/docker"
	"github.com/shubhamku044/containix/internal/ui/views"
)

// MainModel is the main model for the application
type MainModel struct {
	containerList views.ContainerListModel
	logView       views.LogViewModel
	statsView     views.StatsViewModel
	focusLeft     bool
	width         int
	height        int
}

// NewMainModel creates a new main model
func NewMainModel() tea.Model {
	// Create docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		panic(err)
	}

	// Fix: NewContainerListModel doesn't take arguments
	containerList, err := views.NewContainerListModel()
	if err != nil {
		panic(err)
	}

	return MainModel{
		containerList: containerList,
		logView:       views.NewLogViewModel(),
		statsView:     views.NewStatsView(dockerClient),
		focusLeft:     true,
	}
}

// Init initializes the model
func (m MainModel) Init() tea.Cmd {
	return m.containerList.Init()
}

// Update updates the model
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Store full terminal dimensions
		m.width = msg.Width
		m.height = msg.Height

		// Use full height but divide it proportionally
		containerListHeight := m.height * 6 / 10
		statsHeight := m.height - containerListHeight

		// Pass full width/2 to each side component
		containerListMsg := tea.WindowSizeMsg{
			Width:  msg.Width / 2,
			Height: containerListHeight,
		}
		containerListModel, cmd := m.containerList.Update(containerListMsg)
		m.containerList = containerListModel.(views.ContainerListModel)
		cmds = append(cmds, cmd)

		// Log view gets full height on right side
		logMsg := tea.WindowSizeMsg{
			Width:  m.width / 2,
			Height: m.height,
		}
		m.logView, cmd = m.logView.Update(logMsg)
		cmds = append(cmds, cmd)

		// Stats view gets remaining height on left side
		statsMsg := tea.WindowSizeMsg{
			Width:  m.width / 2,
			Height: statsHeight,
		}
		m.statsView, cmd = m.statsView.Update(statsMsg)
		cmds = append(cmds, cmd)

	case views.SelectedContainerMsg:
		// When a container is selected, update the stats view
		cmd := m.statsView.SetContainerID(msg.ID)
		cmds = append(cmds, cmd)

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
		containerListModel, cmd := m.containerList.Update(msg)
		m.containerList = containerListModel.(views.ContainerListModel)
		cmds = append(cmds, cmd)
	} else {
		// Fix: use a temporary variable for the result
		var logViewModel views.LogViewModel
		var cmd tea.Cmd
		logViewModel, cmd = m.logView.Update(msg)
		m.logView = logViewModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the model
func (m MainModel) View() string {
	// Create left side with container list
	containerListView := m.containerList.View()

	// Right side: Split vertical space between statsView and logView
	statsHeight := m.height / 2
	logsHeight := m.height - statsHeight

	// Ensure statsView and logView fit within their allocated heights
	statsView := lipgloss.NewStyle().
		Height(statsHeight).
		Width(m.width / 2).
		Render(m.statsView.View())

	logsView := lipgloss.NewStyle().
		Height(logsHeight).
		Width(m.width / 2).
		Render(m.logView.View())

	// Stack statsView and logView vertically
	rightSide := lipgloss.JoinVertical(lipgloss.Left, statsView, logsView)

	// Calculate exact widths for left and right sections
	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth

	// Apply styles based on focus
	leftStyle := lipgloss.NewStyle().Width(leftWidth).Height(m.height)
	rightStyle := lipgloss.NewStyle().Width(rightWidth).Height(m.height)

	if m.focusLeft {
		leftStyle = leftStyle.BorderForeground(lipgloss.Color("62")).
			BorderStyle(lipgloss.RoundedBorder())
	} else {
		rightStyle = rightStyle.BorderForeground(lipgloss.Color("62")).
			BorderStyle(lipgloss.RoundedBorder())
	}

	// Render styled views
	styledLeft := leftStyle.Render(containerListView)
	styledRight := rightStyle.Render(rightSide)

	// Join horizontally with proper dimensions
	return lipgloss.JoinHorizontal(lipgloss.Top, styledLeft, styledRight)
}
