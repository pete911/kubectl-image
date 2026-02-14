package main

import (
	"os"

	"github.com/pete911/kubectl-image/cmd"
)

var Version = "dev"

func main() {

	cmd.Version = Version
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
