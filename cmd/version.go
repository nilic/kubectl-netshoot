package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print kubectl-netshoot version",
	Long:  `Print kubectl-netshoot version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kubectl-netshoot v%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
