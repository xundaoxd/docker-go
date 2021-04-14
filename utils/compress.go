package utils

import (
	"log"
	"os/exec"
)

func DeCompressTarGz(tar, root string) {
	cmd := exec.Command("tar", "xvf", tar, "-C", root)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
