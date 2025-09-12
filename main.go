package main

import (
	"math/rand"
	"time"

	"github.com/neuralmux/devpod-provider-upcloud/cmd"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	cmd.Execute()
}