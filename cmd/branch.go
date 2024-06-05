package cmd

import (
	"fmt"
	"os"

	"github.com/jplindgren/bastard-git/internal/repository"
	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List, create, or delete branches",
	Run: func(cmd *cobra.Command, args []string) {
		repo := repository.GetRepository(cfg.user)
		if !repo.IsValid() {
			os.Exit(1)

		}

		branch, err := repo.GetCurrentBranch()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout, branch)
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
