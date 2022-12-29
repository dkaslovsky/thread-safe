package main

import (
	"fmt"
	"os"

	"github.com/dkaslovsky/thread-safe/cmd"
)

// TODO: name vs path everywhere

const (
	name    = "thread-safe"
	version = "0.0.1" // hardcode version for now
)

func main() {
	err := cmd.Run(name, version, os.Args)
	if err != nil {
		fmt.Printf("%s: %v\n", name, err)
		os.Exit(1)
	}
}
