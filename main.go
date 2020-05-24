package main

import (
	"fmt"

	"github.com/docker/distribution/context"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
)

func main() {
	fmt.Println("Helow")
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.30", nil, nil)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}
}
