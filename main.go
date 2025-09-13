package main

import (
	"github.com/neuralmux/devpod-provider-upcloud/cmd"
)

func main() {
	// As of Go 1.20, the random number generator is automatically seeded
	cmd.Execute()
}
