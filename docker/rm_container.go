package docker

import (
	"os"
	"path"
)

func rmContainer(container_id string) {
	if err := os.RemoveAll(path.Join(docker_cfg.ContainerRoot, container_id)); err != nil {
		panic(err.Error())
	}
}

func RemoveContainerAction(args []string) {
	for _, a := range args {
		rmContainer(a)
	}
}
