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
	pflag.StringVar(&imageTag, "image-tag", "", "")
	pflag.Parse()

	// form full image name from base and tag
	if imageTag != "" {
		fullImageName = baseImageName + ":" + imageTag
	} else {
		fullImageName = baseImageName + ":latest"
	}

	// set image for child commands
	for _, c := range cmd.GetRootCmd().Commands() {
		if c.Name() == "run" || c.Name() == "debug" {
			c.Flags().Set("image", fullImageName)
		}
	}

	// execute root command
	cmd.Execute()
}
