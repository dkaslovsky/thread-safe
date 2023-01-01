package main

import (
	"fmt"
	"os"

	"github.com/dkaslovsky/thread-safe/cmd"
)

const (
	name = "thread-safe"
)

var version string // set by build ldflags

func main() {
	err := cmd.Run(name, version)
	if err != nil {
		fmt.Printf("%s: %v\n", name, err)
		os.Exit(1)
	}
}
