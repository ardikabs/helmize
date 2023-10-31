package main

import (
	"os"

	"github.com/ardikabs/kasque/cmd"
)

func main() {
	root := cmd.MakeRoot()
	root.AddCommand(cmd.MakeVersion())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
