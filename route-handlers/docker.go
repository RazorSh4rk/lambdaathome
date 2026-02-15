package api

import (
	"github.com/RazorSh4rk/lambdaathome/types"
	commands "github.com/RazorSh4rk/lambdaathome/docker-commands"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
)

type dockerClient interface {
	Kill(ID string)
	RemoveContainer(ID string)
	RemoveImage(tag string)
	RunDetached(lambda types.LambdaFun) string
	ListRunning() []dockerTypes.Container
	ListInstalledImages() []image.Summary
	BuildImage(function types.LambdaFun)
	Close()
}

var newDockerClient = func() (dockerClient, error) {
	return commands.NewClient()
}
