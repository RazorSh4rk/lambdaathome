package commands

import (
	"context"
	"log"
	"strings"

	"github.com/RazorSh4rk/f"
	t "github.com/RazorSh4rk/lambdaathome/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
)

func (client Client) RunDetached(lambda t.LambdaFun) string {
	log.Println("Running image", lambda.Name)
	ctx := context.Background()

	installed := client.ListInstalledImages()
	downloaded := f.From(installed).Has(func(i image.Summary) bool {
		return f.From(i.RepoTags).Has(func(tag string) bool {
			return strings.HasPrefix(tag, lambda.Name)
		})
	})

	if downloaded {
		log.Println("Image already downloaded")
	} else {
		log.Println("Downloading image")
		client.PullImage(lambda.Name)
	}

	ID := client.CreateContainer(lambda.Name)

	if err := client.c.ContainerStart(ctx, ID, container.StartOptions{}); err != nil {
		log.Fatal(err)
	}

	log.Printf("%s started as %s", lambda.Name, ID)
	return ID
}

func (client Client) PullImage(imageName string) {
	ctx := context.Background()
	out, err := client.c.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
}

func (client Client) CreateContainer(imageName string) string {
	ctx := context.Background()

	resp, err := client.c.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")
	if err != nil {
		log.Fatal(err)
	}

	return resp.ID
}
