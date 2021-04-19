package docker

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/xundaoxd/docker-go/utils"
)

func checkImage(fpath string, d fs.DirEntry, err error) error {
	m_path := path.Join(fpath, "manifests.json")
	if utils.IsExist(m_path) {
		img_path := fpath[len(docker_cfg.DockerRoot)+1:]
		idx := strings.LastIndex(img_path, "-")
		fmt.Printf("%s:%s\n", img_path[:idx], img_path[idx+1:])
		return fs.SkipDir
	}
	return nil
}

func ListImageAction(args []string) {
	if err := filepath.WalkDir(docker_cfg.DockerRoot, checkImage); err != nil {
		panic(err.Error())
	}
}
