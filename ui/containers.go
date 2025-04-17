package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type containerListModel struct {
	list list.Model
}

func newContainerList() containerListModel {
	items := []list.Item{
		listItem("dev-db"),
		listItem("api-server"),
		listItem("redis"),
	}
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 30, 20)
	l.Title = "Containers"
	return containerListModel{list: l}
}

func (m containerListModel) Update(msg tea.Msg) (containerListModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m containerListModel) View() string {
	return lipgloss.NewStyle().Padding(1).Width(30).Render(m.list.View())
}

type listItem string

func (i listItem) Title() string       { return string(i) }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return string(i) }
