package commands

import (
	"context"
	"fmt"
	"io"
	"log"

	t "github.com/RazorSh4rk/lambdaathome/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
)

func (client Client) BuildImage(function t.LambdaFun) {
	ctx := context.Background()
	dockerfilePath := fmt.Sprintf(function.Source)

	tar, err := archive.TarWithOptions(dockerfilePath, &archive.TarOptions{})
	if err != nil {
		panic(err)
	}

	res, err := client.c.ImageBuild(ctx, tar, types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{function.Tag},
	})
	if err != nil {
		panic(err)
	}

	message, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
	log.Println(string(message))
}
