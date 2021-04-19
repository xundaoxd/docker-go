package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/xundaoxd/docker-go/net"
	"github.com/xundaoxd/docker-go/utils"
)

func pullImage(arg string) {
	registry := docker_cfg.RegistryUrl
	var img_name string = arg
	var img_tag string = "latest"
	if idx := strings.Index(arg, ":"); idx != -1 {
		img_name = arg[:idx]
		img_tag = arg[idx+1:]
	}
	img_root := path.Join(docker_cfg.DockerRoot, img_name+"-"+img_tag)
	manifests_url := fmt.Sprintf("%s/%s/%s/%s", registry, img_name, "manifests", img_tag)
	manifests_path := path.Join(img_root, "manifests.json")
	if utils.IsExist(manifests_path) {
		return
	}
	if err := os.MkdirAll(img_root, 0755); err != nil {
		panic(err.Error())
	}

	var manifests_json map[string]interface{}
	if buf, err := net.DownloadFile(manifests_url, manifests_path, 0644); err != nil {
		panic(err.Error())
	} else {
		if err := json.Unmarshal(buf, &manifests_json); err != nil {
			panic(err.Error())
		}
	}

	if manifests_json["errors"] != nil {
		panic(manifests_json["errors"])
	}
	for _, layer := range manifests_json["fsLayers"].([]interface{}) {
		blobSum := layer.(map[string]interface{})["blobSum"].(string)
		hash := strings.Split(blobSum, ":")
		layer_url := fmt.Sprintf("%s/%s/%s/%s:%s", registry, img_name, "blobs", hash[0], hash[1])
		layer_root := path.Join(img_root, hash[1])
		layer_file := path.Join(layer_root, "layer.tar.gz")
		if err := os.MkdirAll(layer_root, 0755); err != nil {
			panic(err.Error())
		}
		_, err := net.DownloadFile(layer_url, layer_file, 0644)
		if err != nil {
			panic(err.Error())
		}
		if err := os.MkdirAll(path.Join(img_root, "root"), 0755); err != nil {
			panic(err.Error())
		}
		utils.DeCompressTarGz(layer_file, path.Join(img_root, "root"))
	}
}

func PullAction(args []string) {
	for _, arg := range args {
		pullImage(arg)
	}
}
