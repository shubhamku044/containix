package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shubhamku044/containix/internal/ui/views"
)

// MainModel is the main model for the application
type MainModel struct {
	containerList views.ContainerListModel
	logView       views.LogViewModel
	focusLeft     bool
}

// NewMainModel creates a new main model
func NewMainModel() tea.Model {
	containerList, err := views.NewContainerListModel()
	if err != nil {
		panic(err)
	}
	
	return MainModel{
		containerList: containerList,
		logView:       views.NewLogViewModel(),
		focusLeft:     true,
	}
}

// Init initializes the model
func (m MainModel) Init() tea.Cmd {
	return m.containerList.Init()
}

// Update updates the model
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		newModel, cmd = m.containerList.Update(msg)
		m.containerList = newModel.(views.ContainerListModel)
	}

	return m, cmd
}

// View renders the model
func (m MainModel) View() string {
	left := m.containerList.View()
	right := m.logView.View()
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}