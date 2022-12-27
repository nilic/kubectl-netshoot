package cmd

import (
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
	hostNetwork bool
	imageTag    string

	rootCmd = &cobra.Command{
		Use:   "kubectl-netshoot",
		Short: "kubectl plugin for spinning up netshoot container for network troubleshooting.",
		Long:  "kubectl plugin for spinning up netshoot container for network troubleshooting.",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&hostNetwork,
		"host-network", false, "(applicable to \"run\" only) whether to spin up netshoot on the host's network namespace")
	rootCmd.PersistentFlags().StringVar(&imageTag,
		"image-tag", "latest", "netshoot container image tag to use")
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDiscoveryBurst(350).WithDiscoveryQPS(50.0)
	matchVersionKubeConfigFlags := kcmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := kcmdutil.NewFactory(matchVersionKubeConfigFlags)
	ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	debugCmd := debug.NewCmdDebug(f, ioStreams)
	debugCmd.SetHelpTemplate("TODO: debug help")
	debugCmd.Flags().Set("stdin", "true")
	debugCmd.Flags().Set("tty", "true")
	rootCmd.AddCommand(debugCmd)

	runCmd := run.NewCmdRun(f, ioStreams)
	runCmd.SetHelpTemplate("TODO: run help")
	runCmd.Flags().Set("stdin", "true")
	runCmd.Flags().Set("tty", "true")
	runCmd.Flags().Set("rm", "true")
	rootCmd.AddCommand(runCmd)

	kubeConfigFlags.AddFlags(rootCmd.PersistentFlags())
	matchVersionKubeConfigFlags.AddFlags(rootCmd.PersistentFlags())
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}
