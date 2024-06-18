package commands

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (client Client) ListRunning() []types.Container {
	ctx := context.Background()

	containers, err := client.c.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		panic(err)
	}

	return containers
}
