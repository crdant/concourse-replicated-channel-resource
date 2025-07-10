package main

import (
	"os"
	"path/filepath"

	"github.com/replicatedhq/concourse-replicated-channel-resource/internal/resource"
)

func main() {
	command := filepath.Base(os.Args[0])

	switch command {
	case "check":
		resource.Check()
	case "in":
		resource.In()
	case "out":
		resource.Out()
	default:
		panic("unknown command: " + command)
	}
}
