package docker

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/xundaoxd/docker-go/utils"
)

func ListContainerAction(args []string) {
	var quiet bool
	cmdline := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cmdline.BoolVar(&quiet, "q", false, "")
	cmdline.Parse(args)

	dir, err := ioutil.ReadDir(DockerCfg.ContainerRoot)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range dir {
		if !v.IsDir() {
			continue
		}
		cfg_path := path.Join(DockerCfg.ContainerRoot, v.Name(), "config.json")
		if utils.IsExist(cfg_path) {
			buf, err := ioutil.ReadFile(cfg_path)
			if err != nil {
				continue
			}
			var cfg map[string]interface{}
			if err := json.Unmarshal(buf, &cfg); err != nil {
				continue
			}
			var name string
			for k, v := range cfg {
				if k == "Name" {
					name = v.(string)
				}
			}
			if quiet {
				fmt.Printf("%s\n", v.Name())
			} else {
				fmt.Printf("Id: %s, Name: %s\n", v.Name(), name)
			}
		}
	}
}
