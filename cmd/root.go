package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/debug"
	"k8s.io/kubectl/pkg/cmd/run"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
)

var (
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
	debugCmd, err := getCmdWithFlags("debug")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	debugCmd.SetHelpTemplate("TODO: debug help")
	debugCmd.Flags().Set("image", "nicolaka/netshoot")
	debugCmd.Flags().Set("stdin", "true")
	debugCmd.Flags().Set("tty", "true")
	rootCmd.AddCommand(debugCmd)

	runCmd, err := getCmdWithFlags("run")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	runCmd.SetHelpTemplate("TODO: run help")
	runCmd.Flags().Set("image", "nicolaka/netshoot")
	runCmd.Flags().Set("stdin", "true")
	runCmd.Flags().Set("tty", "true")
	runCmd.Flags().Set("rm", "true")
	rootCmd.AddCommand(runCmd)
}

func getCmdWithFlags(cmdName string) (*cobra.Command, error) {
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDiscoveryBurst(350).WithDiscoveryQPS(50.0)
	matchVersionKubeConfigFlags := kcmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := kcmdutil.NewFactory(matchVersionKubeConfigFlags)
	ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	var cmd *cobra.Command
	switch cmdName {
	case "debug":
		cmd = debug.NewCmdDebug(f, ioStreams)
	case "run":
		cmd = run.NewCmdRun(f, ioStreams)
	default:
		return nil, fmt.Errorf("invalid command type, expected (debug, run), got %s", cmdName)
	}

	kubeConfigFlags.AddFlags(cmd.PersistentFlags())
	matchVersionKubeConfigFlags.AddFlags(cmd.PersistentFlags())
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	return cmd, nil
}
