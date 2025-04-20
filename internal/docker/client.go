package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func FetchContainers() []types.Container {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil // Return nil on error
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil // Return nil on error
	}

	return containers
}

func FetchImages() []types.ImageSummary {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil // Return nil on error
	}
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return nil
	}

	return images
}
