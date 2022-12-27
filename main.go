package main

import (
	"github.com/nilic/kubectl-netshoot/cmd"
	"github.com/spf13/pflag"

	// Import to initialize client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	imageTag      string
	fullImageName string
)

const baseImageName = "nicolaka/netshoot"

func main() {
	pflag.StringVarP(&imageTag, "image-tag", "t", "latest", "netshoot container image tag to use")
	pflag.Parse()

	fullImageName = baseImageName + ":" + imageTag

	for _, c := range cmd.GetRootCmd().Commands() {
		if c.Name() == "run" || c.Name() == "debug" {
			c.Flags().Set("image", fullImageName)
		}
	}

	cmd.Execute()
}
