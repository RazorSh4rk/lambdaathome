package commands

import (
	"context"
	"log"
)

func (client Client) Kill(ID string) {
	ctx := context.Background()
	if err := client.c.ContainerKill(ctx, ID, "SIGKILL"); err != nil {
		log.Fatal(err)
	}
}
