package main

import (
	"os"

	"github.com/shotah/forgecli/forgecli"
)

func main() {
	os.Exit(forgecli.CLI(os.Args[1:]))
}
