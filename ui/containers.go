package ui

import (
	"context"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type logsModel struct {
	viewport    viewport.Model
	content     string
	width       int
	height      int
	parentModel tea.Model
}

func (m logsModel) Init() tea.Cmd {
	return nil
}

func (m logsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m logsModel) View() string {
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

type containerListModel struct {
	list         list.Model
	dockerClient *client.Client
	err          error
	logs         string
	showingLogs  bool
	width        int
	height       int
	logsViewport viewport.Model
}

type containersFetchedMsg struct {
	items []list.Item
}
type errMsg struct {
	err error
}

func newContainerList() (containerListModel, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.48"))
	if err != nil {
		return containerListModel{}, err
	}

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Containers"
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(2)

	vp := viewport.New(0, 0)

	return containerListModel{
		list:         l,
		dockerClient: cli,
		logsViewport: vp,
	}, nil
}

func (m containerListModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchContainers(),
		tea.EnterAltScreen,
	)
}

func (m *containerListModel) handleWindowSize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height

	listWidth := msg.Width - 4
	listHeight := msg.Height - 6

	m.list.SetSize(listWidth, listHeight)

	m.logsViewport.Width = msg.Width - 4
	m.logsViewport.Height = msg.Height - 2
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
				id:     c.ID,
				title:  name,
				status: c.State,
			}
		}
		return containersFetchedMsg{items}
	}
}

func (m containerListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-4, msg.Height-6)

	case containersFetchedMsg:
		m.list.SetItems(msg.items)
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return m, m.fetchContainers()
		case "s":
			if selectedItem, ok := m.list.SelectedItem().(listItem); ok {
				return m, tea.Sequence(
					m.stopContainer(selectedItem.id),
					m.fetchContainers(),
				)
			}
		case "t":
			if selectedItem, ok := m.list.SelectedItem().(listItem); ok {
				return m, tea.Sequence(
					m.startContainer(selectedItem.id),
					m.fetchContainers(),
				)
			}
		case "x":
			if selectedItem, ok := m.list.SelectedItem().(listItem); ok {
				return m, tea.Sequence(
					m.restartContainer(selectedItem.id),
					m.fetchContainers(),
				)
			}
		case "l":
			if selectedItem, ok := m.list.SelectedItem().(listItem); ok {
				m.showingLogs = true
				return m, m.fetchLogs(selectedItem.id)
			}
		case "q":
			if m.showingLogs {
				m.showingLogs = false
				m.logs = ""
			} else {
				return m, tea.Quit
			}
		}

	case logsMsg:
		m.logs = msg.logs
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *containerListModel) stopContainer(containerID string) tea.Cmd {
	return func() tea.Msg {
		err := m.dockerClient.ContainerStop(context.Background(), containerID, container.StopOptions{})
		if err != nil {
			return errMsg{err}
		}
		return nil
	}
}

func (m *containerListModel) startContainer(containerID string) tea.Cmd {
	return func() tea.Msg {
		err := m.dockerClient.ContainerStart(context.Background(), containerID, container.StartOptions{})
		if err != nil {
			return errMsg{err}
		}
		return nil
	}
}

func (m *containerListModel) restartContainer(containerID string) tea.Cmd {
	return func() tea.Msg {
		err := m.dockerClient.ContainerRestart(context.Background(), containerID, container.StopOptions{})
		if err != nil {
			return errMsg{err}
		}
		return nil
	}
}

func (m *containerListModel) fetchLogs(containerID string) tea.Cmd {
	return func() tea.Msg {
		reader, err := m.dockerClient.ContainerLogs(context.Background(), containerID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
		if err != nil {
			return errMsg{err}
		}
		defer reader.Close()

		_, err = io.ReadAll(reader)
		if err != nil {
			return errMsg{err}
		}
		return nil
	}
}

type logsMsg struct {
	logs string
}

func (m containerListModel) View() string {
	if m.err != nil {
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Padding(1, 2).
			Render("Error: " + m.err.Error() + "\nPress R to retry")
	}

	if m.showingLogs {
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Padding(1, 2).
			Render(m.logsViewport.View() + "\n\nPress 'q' to return")
	}

	listView := lipgloss.NewStyle().
		Width(m.width - 4).
		Height(m.height - 6).
		Render(m.list.View())

	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("\n  s: stop • t: start • x: restart • l: logs • r: refresh • q: quit")

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Padding(1, 2).
		Render(listView + helpText)
}

type listItem struct {
	id     string
	title  string
	status string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.status }
func (i listItem) FilterValue() string { return i.title }
