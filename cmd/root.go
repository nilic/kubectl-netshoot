package cmd

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/cmd/debug"
	"k8s.io/kubectl/pkg/cmd/run"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
)

const (
	compTimeout         = 2 * time.Second
	hostNetworkOverride = "{\"spec\": {\"hostNetwork\": true}}"
)

var (
	hostNetwork bool
	imageName   string
	imageTag    string

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
	rootCmd.PersistentFlags().StringVar(&imageName,
		"image-name", "nicolaka/netshoot", "netshoot container image to use")
	rootCmd.PersistentFlags().StringVar(&imageTag,
		"image-tag", "latest", "netshoot container image tag to use")

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDiscoveryBurst(350).WithDiscoveryQPS(50.0)
	matchVersionKubeConfigFlags := kcmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := kcmdutil.NewFactory(matchVersionKubeConfigFlags)
	ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	compFuncPods := func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := getKubernetesClient(f)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		namespace, err := getActiveNamespace(cmd, f)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		candidatePodNames, err := getCandidatePodNames(client, namespace)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return candidatePodNames, cobra.ShellCompDirectiveDefault
	}

	debugCmd := debug.NewCmdDebug(f, ioStreams)
	debugCmd.SetHelpTemplate(debugHelp)
	debugCmd.Short = debugShort
	debugCmd.Flags().Set("stdin", "true")
	debugCmd.Flags().Set("tty", "true")
	debugCmd.ValidArgsFunction = compFuncPods
	rootCmd.AddCommand(debugCmd)

	runCmd := run.NewCmdRun(f, ioStreams)
	runCmd.SetHelpTemplate(runHelp)
	runCmd.Short = runShort
	runCmd.Flags().Set("stdin", "true")
	runCmd.Flags().Set("tty", "true")
	runCmd.Flags().Set("rm", "true")
	runCmd.ValidArgsFunction = compFuncPods
	rootCmd.AddCommand(runCmd)

	kubeConfigFlags.AddFlags(rootCmd.PersistentFlags())
	matchVersionKubeConfigFlags.AddFlags(rootCmd.PersistentFlags())
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	rootCmd.RegisterFlagCompletionFunc("namespace", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := getKubernetesClient(f)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		candidateNamespaces, err := getCandidateNamespaces(client)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return candidateNamespaces, cobra.ShellCompDirectiveDefault
	})
}

func setFlagsForChildCmds(cmd *cobra.Command) {
	fullImageName := imageName + ":" + imageTag

	if cmd.Name() == "run" || cmd.Name() == "debug" {
		cmd.Flags().Set("image", fullImageName)
	}

	if cmd.Name() == "run" && hostNetwork {
		cmd.Flags().Set("overrides", hostNetworkOverride)
	}

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

func getActiveNamespace(cmd *cobra.Command, f kcmdutil.Factory) (string, error) {
	flagNamespace := cmd.Flag("namespace")
	if flagNamespace != nil && flagNamespace.Value.String() != "" {
		return flagNamespace.Value.String(), nil
	}

	namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return "", err
	}

	return namespace, nil
}

func getCandidatePodNames(client *kubernetes.Clientset, namespace string) ([]string, error) {
	ctx, ctxCancelFunc := context.WithTimeout(context.Background(), compTimeout)
	defer ctxCancelFunc()

	pods, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}

	return podNames, nil
}

func getCandidateNamespaces(client *kubernetes.Clientset) ([]string, error) {
	ctx, ctxCancelFunc := context.WithTimeout(context.Background(), compTimeout)
	defer ctxCancelFunc()

	nsList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nsNames []string
	for _, ns := range nsList.Items {
		nsNames = append(nsNames, ns.Name)
	}

	return nsNames, nil
}

func getKubernetesClient(f kcmdutil.Factory) (*kubernetes.Clientset, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
