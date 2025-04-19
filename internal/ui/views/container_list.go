package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/shubhamku044/containix/internal/docker"
)

// ContainerListModel represents the container list view
type ContainerListModel struct {
	list         list.Model
	dockerClient *docker.Client
	err          error
	logs         string
	showingLogs  bool
	width        int
	height       int
	logsViewport viewport.Model
}

// ContainerItem represents a container in the list
type ContainerItem struct {
	id     string
	title  string
	status string
}

func (i ContainerItem) Title() string       { return i.title }
func (i ContainerItem) Description() string { return i.status }
func (i ContainerItem) FilterValue() string { return i.title }

// ContainersFetchedMsg is sent when containers are fetched
type ContainersFetchedMsg struct {
	Items []list.Item
}

// ErrMsg is sent when an error occurs
type ErrMsg struct {
	Err error
}

// LogsMsg is sent when logs are fetched
type LogsMsg struct {
	Logs string
}

// NewContainerListModel creates a new container list model
func NewContainerListModel() (ContainerListModel, error) {
	cli, err := docker.NewClient()
	if err != nil {
		return ContainerListModel{}, err
	}

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Containers"
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(2)

	vp := viewport.New(0, 0)

	return ContainerListModel{
		list:         l,
		dockerClient: cli,
		logsViewport: vp,
	}, nil
}

// Init initializes the model
func (m ContainerListModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchContainers(),
		tea.EnterAltScreen,
	)
}

// fetchContainers fetches containers from Docker
func (m *ContainerListModel) fetchContainers() tea.Cmd {
	return func() tea.Msg {
		containers, err := m.dockerClient.ListContainers()
		if err != nil {
			return ErrMsg{Err: err}
		}

		items := make([]list.Item, len(containers))
		for i, c := range containers {
			items[i] = ContainerItem{
				id:     c.ID,
				title:  c.Name,
				status: c.Status,
			}
		}
		return ContainersFetchedMsg{Items: items}
	}
}

// stopContainer stops a container
func (m *ContainerListModel) stopContainer(containerID string) tea.Cmd {
	return func() tea.Msg {
		err := m.dockerClient.StopContainer(containerID)
		if err != nil {
			return ErrMsg{Err: err}
		}
		return nil
	}
}

// startContainer starts a container
func (m *ContainerListModel) startContainer(containerID string) tea.Cmd {
	return func() tea.Msg {
		err := m.dockerClient.StartContainer(containerID)
		if err != nil {
			return ErrMsg{Err: err}
		}
		return nil
	}
}

// restartContainer restarts a container
func (m *ContainerListModel) restartContainer(containerID string) tea.Cmd {
	return func() tea.Msg {
		err := m.dockerClient.RestartContainer(containerID)
		if err != nil {
			return ErrMsg{Err: err}
		}
		return nil
	}
}

// fetchLogs fetches logs for a container
func (m *ContainerListModel) fetchLogs(containerID string) tea.Cmd {
	return func() tea.Msg {
		logs, err := m.dockerClient.GetContainerLogs(containerID)
		if err != nil {
			return ErrMsg{Err: err}
		}
		return LogsMsg{Logs: logs}
	}
}

// Update updates the model
func (m ContainerListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-4, msg.Height-6)

	case ContainersFetchedMsg:
		m.list.SetItems(msg.Items)
		return m, nil

	case ErrMsg:
		m.err = msg.Err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return m, m.fetchContainers()
		case "s":
			if selectedItem, ok := m.list.SelectedItem().(ContainerItem); ok {
				return m, tea.Sequence(
					m.stopContainer(selectedItem.id),
					m.fetchContainers(),
				)
			}
		case "t":
			if selectedItem, ok := m.list.SelectedItem().(ContainerItem); ok {
				return m, tea.Sequence(
					m.startContainer(selectedItem.id),
					m.fetchContainers(),
				)
			}
		case "x":
			if selectedItem, ok := m.list.SelectedItem().(ContainerItem); ok {
				return m, tea.Sequence(
					m.restartContainer(selectedItem.id),
					m.fetchContainers(),
				)
			}
		case "l":
			if selectedItem, ok := m.list.SelectedItem().(ContainerItem); ok {
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

	case LogsMsg:
		m.logs = msg.Logs
		m.logsViewport.SetContent(msg.Logs)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the model
func (m ContainerListModel) View() string {
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