package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "return the current bgit version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stdout, "0.0.5")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
