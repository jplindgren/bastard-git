package cmd

import (
	"fmt"
	"os"

	"github.com/jplindgren/bastard-git/internal/repository"
	"github.com/spf13/cobra"
)

var createBranchCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new branch",
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("invalid arguments provided\nprovide a branch name")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		repo := repository.GetRepository(cfg.user)
		if !repo.IsValid() {
			os.Exit(1)
		}
		err := repo.CreateBranch(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "New branch created.\nSwitched to a new branch %s"+args[0])
	},
}

func init() {
	branchCmd.AddCommand(createBranchCmd)
}
