package commands

import (
	"context"

	"github.com/docker/docker/api/types/image"
)

func (client Client) ListInstalledImages() []image.Summary {
	ctx := context.Background()

	images, err := client.c.ImageList(ctx, image.ListOptions{})
	if err != nil {
		panic(err)
	}

	return images
}
