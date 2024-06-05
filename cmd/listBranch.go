package cmd

import (
	"fmt"
	"os"

	"github.com/jplindgren/bastard-git/internal/repository"
	"github.com/spf13/cobra"
)

var listBranchCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "Show all branches available",
	Run: func(cmd *cobra.Command, args []string) {
		repo := repository.GetRepository(cfg.user)
		if !repo.IsValid() {
			os.Exit(1)
		}

		branches, err := repo.BranchList()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		currBranch, err := repo.GetCurrentBranch()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		var bList string
		for _, branch := range branches {
			if branch == currBranch {
				branch = "* " + branch
			}
			bList = fmt.Sprintf("%s%s\n", bList, branch)
		}

		fmt.Fprintln(os.Stdout, bList)
	},
}

func init() {
	branchCmd.AddCommand(listBranchCmd)
}
