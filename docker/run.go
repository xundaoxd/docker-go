package docker

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/xundaoxd/docker-go/utils"
)

type Volumes []string

func (v *Volumes) Set(val string) error {
	*v = append(*v, val)
	return nil
}

func (v *Volumes) String() string {
	return strings.Join(*v, ";")
}

func RunAction(args []string) {
	var volumes Volumes
	cmdline := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cmdline.Var(&volumes, "v", "volumes to mount")
	cmdline.Parse(args)
	args = cmdline.Args()
	if len(args) < 2 {
		log.Fatal("container run error")

	}
	for idx, v := range volumes {
		v_info := strings.Split(v, ":")
		if len(v_info) != 2 {
			log.Fatal("mount error: ", v)

		}
		if !strings.HasPrefix(v_info[0], "/") {
			if abs, err := filepath.Abs(v_info[0]); err == nil {
				v_info[0] = abs
			} else {
				log.Fatal(err)
			}
			volumes[idx] = strings.Join(v_info, ":")
		}
	}
	command := args[1:]
	pullImage(args[0])
	var img_name string = args[0]
	var img_tag string = "latest"
	if idx := strings.Index(args[0], ":"); idx != -1 {
		img_name = args[0][:idx]
		img_tag = args[0][idx+1:]
	}
	img_root := path.Join(DockerCfg.DockerRoot, img_name+"-"+img_tag)
	lowerdir := path.Join(img_root, "root")

	container_id := rand.Intn(1000000)
	for utils.IsExist(path.Join(DockerCfg.ContainerRoot, fmt.Sprintf("%06d", container_id))) {
		container_id = rand.Intn(1000000)
	}
	container_root := path.Join(DockerCfg.ContainerRoot, fmt.Sprintf("%06d", container_id))
	container_cfg_path := path.Join(container_root, "config.json")
	upperdir := path.Join(container_root, "upperdir")
	if err := os.MkdirAll(upperdir, 0755); err != nil {
		log.Fatal(err)

	}
	workdir := path.Join(container_root, "workdir")
	if err := os.MkdirAll(workdir, 0755); err != nil {
		log.Fatal(err)

	}
	merged := path.Join(container_root, "merged")
	if err := os.MkdirAll(merged, 0755); err != nil {
		log.Fatal(err)

	}

	var container_cfg Container
	container_cfg.Overlay = make(map[string]string)
	container_cfg.Overlay["lowerdir"] = lowerdir
	container_cfg.Overlay["upperdir"] = upperdir
	container_cfg.Overlay["workdir"] = workdir
	container_cfg.Overlay["merged"] = merged
	container_cfg.Volumes = volumes
	if buf, err := json.Marshal(container_cfg); err == nil {
		if err := ioutil.WriteFile(container_cfg_path, buf, 0644); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	cmd := exec.Command(os.Args[0], append([]string{"exec", fmt.Sprintf("%06d", container_id)}, command...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
