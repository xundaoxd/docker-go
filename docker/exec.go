package docker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	if err := syscall.Chroot(root_dir); err != nil {
		panic(err.Error())
	}
	if err := syscall.Chdir("/"); err != nil {
		panic(err.Error())
	}
	proc_dir := path.Join("/proc")
	if err := os.MkdirAll(proc_dir, 0755); err != nil {
		panic(err.Error())
	}
	if err := syscall.Mount("proc", proc_dir, "proc", 0, ""); err != nil {
		panic(err.Error())
	}
	defer func() {
		if err := syscall.Unmount(proc_dir, 0); err != nil {
			panic(err.Error())
		}
	}()
	cmd := exec.Command("env", "-i", strings.Join(command, " "))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err.Error())
	}
}

func ExecAction(args []string) {
	if os.Args[0] == "/proc/self/exe" {
		runSelf(args)
		return
	}
	if len(args) < 2 {
		panic("container run error")

	}
	container_id := args[0]
	command := args[1:]
	container_root := path.Join(docker_cfg.ContainerRoot, container_id)
	container_cfg_path := path.Join(container_root, "config.json")
	if !utils.IsExist(container_cfg_path) {
		panic("container doesn't exist.")
	}
	var container_cfg Container
	if buf, err := ioutil.ReadFile(container_cfg_path); err != nil {
		panic(err.Error())
	} else {
		if err := json.Unmarshal(buf, &container_cfg); err != nil {
			panic(err.Error())
		}
	}

	lowerdir := container_cfg.Overlay["lowerdir"]
	upperdir := container_cfg.Overlay["upperdir"]
	workdir := container_cfg.Overlay["workdir"]
	merged := container_cfg.Overlay["merged"]
	if err := syscall.Mount("overlay", merged, "overlay", 0, fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerdir, upperdir, workdir)); err != nil {
		panic(err.Error())
	}
	defer func() {
		if err := syscall.Unmount(merged, 0); err != nil {
			panic(err.Error())
		}
	}()
	for _, v := range container_cfg.Volumes {
		v_info := strings.Split(v, ":")
		if len(v_info) != 2 {
			panic("mount error: " + v)
		}
		target := path.Join(merged, v_info[1])
		if err := os.MkdirAll(target, 0755); err != nil {
			panic(err.Error())
		}
		if err := syscall.Mount(v_info[0], target, "", syscall.MS_BIND, ""); err != nil {
			panic(err.Error())
		}
		defer func() {
			if err := syscall.Unmount(target, 0); err != nil {
				panic(err.Error())
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
		panic(err.Error())
	}
}
