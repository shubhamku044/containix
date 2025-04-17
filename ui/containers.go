package ui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type containerListModel struct {
	list         list.Model
	dockerClient *client.Client
	err          error
}

type containersFetchedMsg struct {
	items []list.Item
}
type errMsg struct {
	err error
}

func newContainerList() (containerListModel, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return containerListModel{}, err
	}
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 30, 20)
	l.Title = "Containers"

	return containerListModel{
		list:         l,
		dockerClient: cli,
	}, nil
}

func (m containerListModel) Init() tea.Cmd {
	return m.fetchContainers()
}

func (m *containerListModel) fetchContainers() tea.Cmd {
	return func() tea.Msg {
		containers, err := m.dockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
		if err != nil {
			return errMsg{err}
		}

		items := make([]list.Item, len(containers))
		for i, c := range containers {
			name := "Unnamed"
			if len(c.Names) > 0 {
				name = strings.TrimPrefix(c.Names[0], "/")
			}
			items[i] = listItem{
				title:  name,
				status: c.State,
			}
		}
		return containersFetchedMsg{items}
	}
}

func (m containerListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case containersFetchedMsg:
		m.list.SetItems(msg.items)
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "r" {
			return m, m.fetchContainers()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m containerListModel) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\nPress R to retry"
	}
	return lipgloss.NewStyle().Padding(1).Width(30).Render(m.list.View())
}

type listItem struct {
	title  string
	status string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.status }
func (i listItem) FilterValue() string { return i.title }
