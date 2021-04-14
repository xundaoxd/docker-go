package main

import (
	"log"
	"os"

	"github.com/xundaoxd/docker-go/docker"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must specify a command")
	}
	switch os.Args[1] {
	case "pull":
		docker.PullAction(os.Args[2:])
	case "run":
		docker.RunAction(os.Args[2:])
	case "exec":
		docker.ExecAction(os.Args[2:])
	default:
		log.Fatal("undefined command")
	}
}
