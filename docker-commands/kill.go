package commands

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
)

func (client Client) Kill(ID string) {
	ctx := context.Background()
	if err := client.c.ContainerKill(ctx, ID, "SIGKILL"); err != nil {
		log.Println(err)
		return
	}
}

func (client Client) RemoveContainer(ID string) {
	ctx := context.Background()
	if err := client.c.ContainerRemove(ctx, ID, container.RemoveOptions{}); err != nil {
		log.Println(err)
	}
}

func (client Client) RemoveImage(tag string) {
	ctx := context.Background()
	if _, err := client.c.ImageRemove(ctx, tag, image.RemoveOptions{}); err != nil {
		log.Println(err)
	}
}
