package docker

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/xundaoxd/docker-go/utils"
)

func ListContainerAction(args []string) {
	var quiet bool
	cmdline := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cmdline.BoolVar(&quiet, "q", false, "")
	cmdline.Parse(args)

	dir, err := ioutil.ReadDir(docker_cfg.ContainerRoot)
	if err != nil {
		panic(err.Error())
	}
	for _, v := range dir {
		if !v.IsDir() {
			continue
		}
		cfg_path := path.Join(docker_cfg.ContainerRoot, v.Name(), "config.json")
		if utils.IsExist(cfg_path) {
			var container_cfg Container
			if buf, err := ioutil.ReadFile(cfg_path); err != nil {
				continue
			} else {
				if err := json.Unmarshal(buf, &container_cfg); err != nil {
					continue
				}
			}

			if quiet {
				fmt.Printf("%s\n", v.Name())
			} else {
				fmt.Printf("Id: %s, Name: %s\n", v.Name(), container_cfg.Name)
			}
		}
	}
}
