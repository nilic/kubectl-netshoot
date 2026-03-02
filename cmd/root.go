package cmd

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/debug"
	"k8s.io/kubectl/pkg/cmd/run"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
)

var (
	hostNetwork  bool
	imageName    string
	imageTag     string
	nodeSelector map[string]string

	rootCmd = &cobra.Command{
		Use:   "kubectl-netshoot",
		Short: "kubectl plugin for spinning up netshoot container for network troubleshooting.",
		Long:  "kubectl plugin for spinning up netshoot container for network troubleshooting.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setFlagsForChildCmds(cmd)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&hostNetwork,
		"host-network", false, "(\"run\" command only) spin up netshoot on the node's network namespace")
	rootCmd.PersistentFlags().StringToStringVar(&nodeSelector,
		"node-selector", map[string]string{}, "(\"run\" command only) node labels to use as a node selector for scheduling the netshoot pod (e.g. kubernetes.io/os=linux)")
	rootCmd.PersistentFlags().StringVar(&imageName,
		"image-name", "nicolaka/netshoot", "netshoot container image to use")
	rootCmd.PersistentFlags().StringVar(&imageTag,
		"image-tag", "latest", "netshoot container image tag to use")

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDiscoveryBurst(350).WithDiscoveryQPS(50.0)
	matchVersionKubeConfigFlags := kcmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := kcmdutil.NewFactory(matchVersionKubeConfigFlags)
	ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	debugCmd := debug.NewCmdDebug(f, ioStreams)
	debugCmd.SetHelpTemplate(debugHelp)
	debugCmd.Short = debugShort
	debugCmd.Flags().Set("stdin", "true")
	debugCmd.Flags().Set("tty", "true")
	rootCmd.AddCommand(debugCmd)

	runCmd := run.NewCmdRun(f, ioStreams)
	runCmd.SetHelpTemplate(runHelp)
	runCmd.Short = runShort
	runCmd.Flags().Set("stdin", "true")
	runCmd.Flags().Set("tty", "true")
	runCmd.Flags().Set("rm", "true")
	rootCmd.AddCommand(runCmd)

	kubeConfigFlags.AddFlags(rootCmd.PersistentFlags())
	matchVersionKubeConfigFlags.AddFlags(rootCmd.PersistentFlags())
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}

func setFlagsForChildCmds(cmd *cobra.Command) {
	fullImageName := imageName + ":" + imageTag

	if cmd.Name() == "run" || cmd.Name() == "debug" {
		cmd.Flags().Set("image", fullImageName)
	}

	if cmd.Name() == "run" && (hostNetwork || len(nodeSelector) > 0) {
		type specOverride struct {
			HostNetwork  bool              `json:"hostNetwork,omitempty"`
			NodeSelector map[string]string `json:"nodeSelector,omitempty"`
		}
		type podOverride struct {
			Spec specOverride `json:"spec"`
		}
		override := podOverride{}
		override.Spec.HostNetwork = hostNetwork
		if len(nodeSelector) > 0 {
			override.Spec.NodeSelector = nodeSelector
		}
		data, err := json.Marshal(override)
		if err != nil {
			log.Fatalf("error building overrides: %v\n", err)
		}
		cmd.Flags().Set("overrides", string(data))
	}

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
