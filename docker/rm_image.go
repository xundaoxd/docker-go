package docker

import (
	"os"
	"path"
	"strings"
)

func rmImage(image_id string) {
	image_root := path.Join(docker_cfg.DockerRoot, strings.Join(strings.Split(image_id, ":"), "-"))
	if err := os.RemoveAll(image_root); err != nil {
		panic(err.Error())
	}
}

func RemoveImageAction(args []string) {
	for _, a := range args {
		rmImage(a)
	}
}
