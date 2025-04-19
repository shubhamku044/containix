package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shubhamku044/containix/internal/docker"
)

// Styling constants
var (
	statsBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("111")).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	noSelectionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Italic(true)
)

// StatsViewModel displays detailed stats for a selected container
type StatsViewModel struct {
	viewport     viewport.Model
	width        int
	height       int
	containerID  string
	dockerClient *docker.Client
	stats        *docker.ContainerStats
}

// NewStatsView creates a new stats view component
func NewStatsView(dockerClient *docker.Client) StatsViewModel {
	vp := viewport.New(0, 5) // Height will be adjusted based on window size
	return StatsViewModel{
		viewport:     vp,
		dockerClient: dockerClient,
	}
}

// SetContainerID updates the container ID for which stats should be displayed
func (m *StatsViewModel) SetContainerID(id string) tea.Cmd {
	m.containerID = id
	if id == "" {
		m.stats = nil
		return nil
	}
	return m.fetchStats()
}

// fetchStats retrieves stats for the currently selected container
func (m *StatsViewModel) fetchStats() tea.Cmd {
	return func() tea.Msg {
		if m.containerID == "" {
			return nil
		}

		stats, err := m.dockerClient.GetContainerStats(m.containerID)
		if err != nil {
			return ErrMsg{Err: err}
		}

		return ContainerStatsMsg{Stats: stats}
	}
}

// Update handles UI events and updates the view
func (m StatsViewModel) Update(msg tea.Msg) (StatsViewModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Adjust width and height based on terminal size
		m.width = msg.Width
		m.height = msg.Height

		// Subtract padding and border sizes for the viewport
		m.viewport.Width = m.width - statsBoxStyle.GetHorizontalPadding() - statsBoxStyle.GetHorizontalBorderSize()
		m.viewport.Height = m.height - statsBoxStyle.GetVerticalPadding() - statsBoxStyle.GetVerticalBorderSize() - 2 // Reserve space for the title

		// Update content after resize
		m.updateViewportContent()

	case ContainerStatsMsg:
		m.stats = msg.Stats
		m.updateViewportContent()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// updateViewportContent refreshes the content in the viewport based on current stats
func (m *StatsViewModel) updateViewportContent() {
	if m.containerID == "" || m.stats == nil {
		m.viewport.SetContent(noSelectionStyle.Render("Select a container to view stats"))
		return
	}

	// Format the stats into a multi-column layout
	statsContent := ""

	// Row 1: CPU and Memory
	cpuSection := fmt.Sprintf("%s %s\n%s %s",
		labelStyle.Render("CPU:"),
		valueStyle.Render(fmt.Sprintf("%.2f%%", m.stats.CPUPercentage)),
		labelStyle.Render("Memory:"),
		valueStyle.Render(fmt.Sprintf("%s / %s (%.1f%%)",
			formatBytes(uint64ToInt64(m.stats.MemoryUsage)),
			formatBytes(uint64ToInt64(m.stats.MemoryLimit)),
			m.stats.MemoryPercentage)))

	// Row 2: Network
	networkSection := fmt.Sprintf("%s %s / %s",
		labelStyle.Render("Network:"),
		valueStyle.Render(fmt.Sprintf("↓ %s", formatBytes(uint64ToInt64(m.stats.NetworkRx)))),
		valueStyle.Render(fmt.Sprintf("↑ %s", formatBytes(uint64ToInt64(m.stats.NetworkTx)))))

	// Row 3: Block I/O
	ioSection := fmt.Sprintf("%s %s / %s",
		labelStyle.Render("I/O:"),
		valueStyle.Render(fmt.Sprintf("Read: %s", formatBytes(uint64ToInt64(m.stats.BlockRead)))),
		valueStyle.Render(fmt.Sprintf("Write: %s", formatBytes(uint64ToInt64(m.stats.BlockWrite)))))

	// Row 4: PIDs
	pidsSection := fmt.Sprintf("%s %s",
		labelStyle.Render("PIDs:"),
		valueStyle.Render(fmt.Sprintf("%d", m.stats.PIDs)))

	statsContent = lipgloss.JoinVertical(lipgloss.Left,
		cpuSection,
		networkSection,
		ioSection,
		pidsSection)

	m.viewport.SetContent(statsContent)
}

// View renders the stats view
func (m StatsViewModel) View() string {
	// Render the title and content
	title := titleStyle.Render("Container Stats")
	content := m.viewport.View()

	// Render the stats box with adjusted dimensions
	return statsBoxStyle.
		Width(m.width).
		Height(m.height).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, content))
}

// Helper function to format bytes to human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Helper function to safely convert uint64 to int64
func uint64ToInt64(val uint64) int64 {
	// Handle potential overflow for extremely large values
	if val > 9223372036854775807 {
		return 9223372036854775807 // max int64 value
	}
	return int64(val)
}

// ContainerStatsMsg is sent when container stats are fetched
type ContainerStatsMsg struct {
	Stats *docker.ContainerStats
}
