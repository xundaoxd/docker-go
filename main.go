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
	case "images":
		docker.ListImageAction(os.Args[2:])
	case "ps":
		docker.ListContainerAction(os.Args[2:])
	case "rm":
		docker.RemoveContainerAction(os.Args[2:])
	case "rmi":
		docker.RemoveImageAction(os.Args[2:])
	default:
		log.Fatal("undefined command")
	}
}
