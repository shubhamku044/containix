package docker

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Client wraps the Docker client and provides container operations
type Client struct {
	client *client.Client
}

// NewClient creates a new Docker client
func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return nil, err
	}

	return &Client{
		client: cli,
	}, nil
}

// Container represents a Docker container
type Container struct {
	ID     string
	Name   string
	Status string
}

// ListContainers returns a list of all containers
func (c *Client) ListContainers() ([]Container, error) {
	containers, err := c.client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	result := make([]Container, len(containers))
	for i, container := range containers {
		name := "Unnamed"
		if len(container.Names) > 0 {
			name = strings.TrimPrefix(container.Names[0], "/")
		}
		result[i] = Container{
			ID:     container.ID,
			Name:   name,
			Status: container.State,
		}
	}

	return result, nil
}

// StopContainer stops a container
func (c *Client) StopContainer(containerID string) error {
	return c.client.ContainerStop(context.Background(), containerID, nil)
}

// StartContainer starts a container
func (c *Client) StartContainer(containerID string) error {
	return c.client.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(containerID string) error {
	return c.client.ContainerRestart(context.Background(), containerID, nil)
}

// GetContainerLogs returns the logs of a container
func (c *Client) GetContainerLogs(containerID string) (string, error) {
	reader, err := c.client.ContainerLogs(context.Background(), containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(logs), nil
}

// ContainerStats represents statistics of a container
type ContainerStats struct {
	CPUPercentage    float64
	MemoryUsage      uint64 // Changed from int64 to uint64
	MemoryLimit      uint64 // Changed from int64 to uint64
	MemoryPercentage float64
	NetworkRx        uint64 // Changed from int64 to uint64
	NetworkTx        uint64 // Changed from int64 to uint64
	BlockRead        uint64 // Changed from int64 to uint64
	BlockWrite       uint64 // Changed from int64 to uint64
	PIDs             uint64 // Changed from int to uint64
}

// GetContainerStats returns stats for a specific container
func (c *Client) GetContainerStats(containerID string) (*ContainerStats, error) {
	ctx := context.Background()

	// Get stats with stream=false for a one-time stats fetch
	stats, err := c.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	// Parse the stats JSON response
	var statsJSON types.StatsJSON
	if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
		return nil, err
	}

	// Calculate CPU percentage
	cpuPercentage := calculateCPUPercentage(&statsJSON)

	// Calculate memory info
	memoryUsage := statsJSON.MemoryStats.Usage - statsJSON.MemoryStats.Stats["cache"]
	memoryLimit := statsJSON.MemoryStats.Limit
	var memoryPercentage float64
	if memoryLimit > 0 {
		memoryPercentage = float64(memoryUsage) / float64(memoryLimit) * 100.0
	}

	// Network stats
	var rxBytes, txBytes uint64 // Changed from int64 to uint64
	for _, network := range statsJSON.Networks {
		rxBytes += network.RxBytes
		txBytes += network.TxBytes
	}

	// Block IO stats
	var blockRead, blockWrite uint64 // Changed from int64 to uint64
	for _, blkio := range statsJSON.BlkioStats.IoServiceBytesRecursive {
		if blkio.Op == "Read" {
			blockRead = blkio.Value
		} else if blkio.Op == "Write" {
			blockWrite = blkio.Value
		}
	}

	return &ContainerStats{
		CPUPercentage:    cpuPercentage,
		MemoryUsage:      memoryUsage,
		MemoryLimit:      memoryLimit,
		MemoryPercentage: memoryPercentage,
		NetworkRx:        rxBytes,
		NetworkTx:        txBytes,
		BlockRead:        blockRead,
		BlockWrite:       blockWrite,
		PIDs:             statsJSON.PidsStats.Current,
	}, nil
}

// calculateCPUPercentage calculates the CPU usage percentage based on Docker stats
func calculateCPUPercentage(stats *types.StatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuCount := float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
		if cpuCount > 0 {
			return (cpuDelta / systemDelta) * cpuCount * 100.0
		}
	}
	return 0.0
}
