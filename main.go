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
	help          bool
)

const baseImageName = "nicolaka/netshoot"

func main() {
	pflag.StringVarP(&imageTag, "image-tag", "t", "latest", "netshoot container image tag to use")
	pflag.BoolVarP(&help, "help", "h", false, "help for kubectl-netshoot")
	pflag.Usage = func() {
		cmd.GetRootCmd().Usage()
	}
	pflag.Parse()

	fullImageName = baseImageName + ":" + imageTag

	for _, c := range cmd.GetRootCmd().Commands() {
		if c.Name() == "run" || c.Name() == "debug" {
			c.Flags().Set("image", fullImageName)
		}
	}

	cmd.Execute()
	cmd.GetRootCmd().Usage()
}
