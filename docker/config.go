package docker

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type DockerConfig struct {
	RegistryUrl   string
	DockerRoot    string
	ContainerRoot string
}

type Container struct {
	Overlay map[string]string
	Volumes []string
}

var DockerCfg DockerConfig

func init() {

	DOCKER_ROOT := os.Getenv("DOCKER_ROOT")
	if DOCKER_ROOT == "" {
		DOCKER_ROOT = "~/.docker"
	}
	// init DockerCfg
	buf, err := ioutil.ReadFile(path.Join(DOCKER_ROOT, "config.json"))
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(buf, &DockerCfg); err != nil {
		log.Fatal(err)
	}
	if DockerCfg.DockerRoot == "" {
		DockerCfg.DockerRoot = DOCKER_ROOT
	} else if !strings.HasPrefix(DockerCfg.DockerRoot, "/") {
		DockerCfg.DockerRoot = path.Join(DOCKER_ROOT, DockerCfg.DockerRoot)
	}

	if DockerCfg.ContainerRoot == "" {
		DockerCfg.ContainerRoot = DOCKER_ROOT
	} else if !strings.HasPrefix(DockerCfg.ContainerRoot, "/") {
		DockerCfg.ContainerRoot = path.Join(DOCKER_ROOT, DockerCfg.ContainerRoot)
	}

}
