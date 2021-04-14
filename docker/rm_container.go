package docker

import (
	"log"
	"os"
	"path"
)

func rmContainer(container_id string) {
	if err := os.RemoveAll(path.Join(DockerCfg.ContainerRoot, container_id)); err != nil {
		log.Fatal(err)
	}
}

func RemoveContainerAction(args []string) {
	for _, a := range args {
		rmContainer(a)
	}
}
