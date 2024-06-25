package commands

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"slices"
	"strings"

	t "github.com/RazorSh4rk/lambdaathome/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

func (client Client) RunDetached(lambda t.LambdaFun) string {
	log.Println("Running image ", lambda)
	ctx := context.Background()

	ID := client.CreateContainer(lambda)

	if err := client.c.ContainerStart(ctx, ID, container.StartOptions{}); err != nil {
		log.Println(err)
	}

	log.Printf("%s started as %s", lambda.Name, ID)
	return ID
}

func (client Client) PullImage(imageName string) error {
	log.Println("pulling image ", imageName)
	ctx := context.Background()
	out, err := client.c.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

func (client Client) CreateContainer(lambda t.LambdaFun) string {
	log.Println("creating container from ", lambda)

	if slices.Contains([]string{"80", "8080", "9001"}, lambda.Port) {
		log.Println("80, 8080 and 9001 are reserved ports")
		return ""
	}
	portMap := fmt.Sprintf("%s/tcp", lambda.Port)
	log.Printf("mapping port: %s\n", portMap)

	hasVolume := lambda.Volume != ""
	volumes := make(map[string]string, 0)
	if hasVolume {
		volConfig := strings.Split(lambda.Volume, ",")
		for _, volMap := range volConfig {
			mapping := strings.Split(volMap, ":")
			from := mapping[0]
			to := mapping[1]

			volumes[from] = to
		}
	}

	ctx := context.Background()

	network, err := client.c.NetworkCreate(ctx, "lambdaathome", network.CreateOptions{
		Driver: "bridge",
	})

	if err != nil {
		log.Println(err)
	}

	containerConfig := &container.Config{
		Image: lambda.Tag,
		Healthcheck: &container.HealthConfig{
			Test: []string{"CMD-SHELL", "exit 0"},
		},
		ExposedPorts: nat.PortSet{
			nat.Port(portMap): struct{}{},
		},
		// Volumes: map[string]struct{}{
		// 	"/data": {},
		// },
	}

	if hasVolume {
		contConfVols := make(map[string]struct{})
		for k := range volumes {
			contConfVols[k] = struct{}{}
		}
		containerConfig.Volumes = contConfVols
	}

	//fPath, _ := filepath.Abs("./docs")

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port("8080/tcp"): []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: lambda.Port,
				},
			},
		},
		// Mounts: []mount.Mount{
		// 	{
		// 		Type:   mount.TypeBind,
		// 		Source: fPath,
		// 		Target: "/docs",
		// 	},
		// },

		NetworkMode: container.NetworkMode(network.ID),
	}

	if hasVolume {
		hostConfVols := make([]mount.Mount, 0)
		for k, v := range volumes {
			absPath, err := filepath.Abs(k)
			if err != nil {
				log.Printf("error with port mapping: %s -> %s: %s\n", k, v, err)
				continue
			}
			hostConfVols = append(hostConfVols, mount.Mount{
				Type:   mount.TypeBind,
				Source: absPath,
				Target: v,
			})
			log.Printf("mounted %s as %s", absPath, v)
		}
		hostConfig.Mounts = hostConfVols
	}

	resp, err := client.c.ContainerCreate(ctx,
		containerConfig, hostConfig, nil, nil, "")

	if err != nil {
		log.Fatal(err)
	}

	return resp.ID
}
