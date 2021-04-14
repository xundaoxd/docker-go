package docker

import (
	"log"
	"os"
	"path"
	"strings"
)

func rmImage(image_id string) {
	image_root := path.Join(DockerCfg.DockerRoot, strings.Join(strings.Split(image_id, ":"), "-"))
	if err := os.RemoveAll(image_root); err != nil {
		log.Fatal(err)
	}
}

func RemoveImageAction(args []string) {
	for _, a := range args {
		rmImage(a)
	}
}
