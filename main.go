package main

import (
	"fmt"
	"os"

	"github.com/frank-vpl/gnet-cli/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "gnet: %v\n", err)
		os.Exit(1)
	}
}
