package docker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/xundaoxd/docker-go/utils"
)

func runSelf(args []string) {
	root_dir := args[0]
	command := args[1:]
	proc_dir := path.Join(root_dir, "proc")
	if err := os.MkdirAll(proc_dir, 0755); err != nil {
		log.Fatal(err)
	}
	if err := syscall.Mount("proc", proc_dir, "proc", 0, ""); err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := syscall.Unmount("/proc", 0); err != nil {
			log.Println(err)
		}
	}()
	if err := syscall.Chroot(root_dir); err != nil {
		if err := syscall.Unmount(proc_dir, 0); err != nil {
			log.Fatal(err)
		}
		log.Fatal(err)
	}
	if err := syscall.Chdir("/"); err != nil {
		log.Println(err)
		return
	}
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return
	}
}

func ExecAction(args []string) {
	if os.Args[0] == "/proc/self/exe" {
		runSelf(args)
		return
	}
	if len(args) < 2 {
		log.Fatal("container run error")

	}
	container_id := args[0]
	command := args[1:]
	container_root := path.Join(DockerCfg.ContainerRoot, container_id)
	container_cfg_path := path.Join(container_root, "config.json")
	if !utils.IsExist(container_cfg_path) {
		log.Fatal("container doesn't exist.")
	}
	var container_cfg Container
	buf, err := ioutil.ReadFile(container_cfg_path)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(buf, &container_cfg); err != nil {
		log.Fatal(err)
	}
	lowerdir := container_cfg.Overlay["lowerdir"]
	upperdir := container_cfg.Overlay["upperdir"]
	workdir := container_cfg.Overlay["workdir"]
	merged := container_cfg.Overlay["merged"]
	if err := syscall.Mount("overlay", merged, "overlay", 0, fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerdir, upperdir, workdir)); err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if err := syscall.Unmount(merged, 0); err != nil {
			log.Println(err)
		}
	}()
	for _, v := range container_cfg.Volumes {
		v_info := strings.Split(v, ":")
		if len(v_info) != 2 {
			log.Println("mount error: ", v)
			return
		}
		target := path.Join(merged, v_info[1])
		if err := os.MkdirAll(target, 0755); err != nil {
			log.Println(err)
			return
		}
		if err := syscall.Mount(v_info[0], target, "", syscall.MS_BIND, ""); err != nil {
			log.Println(err)
			return
		}
		defer func() {
			if err := syscall.Unmount(target, 0); err != nil {
				log.Println(err)
			}
		}()
	}
	cmd := exec.Command("/proc/self/exe", append([]string{"exec", merged}, command...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return
	}
}
