package cmd

import (
	"github.com/spf13/cobra"
)

// asyncCmd represents the async command
var asyncCmd = &cobra.Command{
	Use:   "async",
	Short: "Works with the async task queue.",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(asyncCmd)
}
