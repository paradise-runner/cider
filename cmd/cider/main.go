package main

import (
	"os"

	"github.com/paradise-runner/cider/internal/cli"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	cli.SetVersion(version)
	if err := cli.Execute(); err != nil {
		// Check if it's an exit error (like diff returning 1 for differences)
		// These should not print usage
		os.Exit(1)
	}
}
