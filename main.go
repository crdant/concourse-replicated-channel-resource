package main

import (
	"os"
	"path/filepath"
)

func main() {
	command := filepath.Base(os.Args[0])
	
	switch command {
	case "check":
		check()
	case "in":
		in()
	case "out":
		out()
	default:
		panic("unknown command: " + command)
	}
}