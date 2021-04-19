package docker

import (
	"encoding/json"
	"io/ioutil"
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
	Name    string
	Overlay map[string]string
	Volumes []string
}

var docker_cfg DockerConfig

func init() {

	DOCKER_ROOT := os.Getenv("DOCKER_ROOT")
	if DOCKER_ROOT == "" {
		DOCKER_ROOT = "~/.docker"
	}
	// init docker_cfg
	if buf, err := ioutil.ReadFile(path.Join(DOCKER_ROOT, "config.json")); err != nil {
		panic(err.Error())
	} else {
		if err := json.Unmarshal(buf, &docker_cfg); err != nil {
			panic(err.Error())
		}
	}

	if docker_cfg.DockerRoot == "" {
		docker_cfg.DockerRoot = DOCKER_ROOT
	} else if !strings.HasPrefix(docker_cfg.DockerRoot, "/") {
		docker_cfg.DockerRoot = path.Join(DOCKER_ROOT, docker_cfg.DockerRoot)
	}

	if docker_cfg.ContainerRoot == "" {
		docker_cfg.ContainerRoot = DOCKER_ROOT
	} else if !strings.HasPrefix(docker_cfg.ContainerRoot, "/") {
		docker_cfg.ContainerRoot = path.Join(DOCKER_ROOT, docker_cfg.ContainerRoot)
	}

}
