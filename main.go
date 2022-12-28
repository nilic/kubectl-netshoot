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
	hostNetwork   bool
	help          bool
)

const (
	baseImageName       = "nicolaka/netshoot"
	hostNetworkOverride = "{\"spec\": {\"hostNetwork\": true}}"
)

func main() {
	pflag.BoolVar(&hostNetwork, "host-network", false, "(applicable to \"run\" command only) whether to spin up netshoot on the host's network namespace")
	pflag.StringVar(&imageTag, "image-tag", "latest", "netshoot container image tag to use")
	pflag.BoolVarP(&help, "help", "h", false, "help for kubectl-netshoot")
	pflag.Usage = func() {
		cmd.GetRootCmd().Usage()
	}
	pflag.Parse()

	fullImageName = baseImageName + ":" + imageTag

	for _, c := range cmd.GetRootCmd().Commands() {
		if c.Name() != "version" {
			c.Flags().Set("image", fullImageName)
			if c.Name() == "run" && hostNetwork {
				c.Flags().Set("overrides", hostNetworkOverride)
			}
		}
	}

	cmd.Execute()
}
