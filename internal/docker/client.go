package docker

import (
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Client wraps the Docker client and provides container operations
type Client struct {
	client *client.Client
}

// NewClient creates a new Docker client
func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.48"))
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
	containers, err := c.client.ContainerList(context.Background(), container.ListOptions{All: true})
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
	return c.client.ContainerStop(context.Background(), containerID, container.StopOptions{})
}

// StartContainer starts a container
func (c *Client) StartContainer(containerID string) error {
	return c.client.ContainerStart(context.Background(), containerID, container.StartOptions{})
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(containerID string) error {
	return c.client.ContainerRestart(context.Background(), containerID, container.StopOptions{})
}

// GetContainerLogs returns the logs of a container
func (c *Client) GetContainerLogs(containerID string) (string, error) {
	reader, err := c.client.ContainerLogs(context.Background(), containerID, container.LogsOptions{
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