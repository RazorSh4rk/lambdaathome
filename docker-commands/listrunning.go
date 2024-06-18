package commands

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (client Client) ListRunning() []types.Container {
	ctx := context.Background()

	containers, err := client.c.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return containers
}
