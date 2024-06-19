package cmd

import (
	"fmt"
	"os"

	"github.com/jplindgren/bastard-git/internal/repository"
	"github.com/spf13/cobra"
)

// var branchName string
var checkoutCmd = &cobra.Command{
	Use:                "checkout",
	Short:              "Switch/create branches or restore working tree files",
	DisableFlagParsing: true,
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("invalid arguments provided\nprovide a branch name or -b to create a new branch")
		}
		// Run the custom validation logic
		if len(args) == 2 && args[0] != "-b" {
			return fmt.Errorf("invalid argument provided: %s", args[1])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		repo := repository.GetRepository(cfg.user)
		if !repo.IsValid() {
			os.Exit(1)
		}

		if len(args) == 0 {
			fmt.Println("Please provide an argument")
			os.Exit(1)
		}

		if len(args) == 2 && args[0] != "-b" {
			fmt.Println("Please provide -b as argument")
			os.Exit(1)
		}

		if len(args) == 1 { //switch branches
			//get last tree
			_, treeHash, err := repo.GetBranchTip(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if treeHash == "" {
				fmt.Fprintln(os.Stdin, "Branch not found")
				os.Exit(1)
			}

			err = repository.RecreateWorkingTree(treeHash)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			repo.SetHead(args[0])

			os.Stdout.WriteString("Switched to branch " + args[0])
		}

		if len(args) == 2 {
			err := repo.CreateBranch(args[1])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Fprintf(os.Stdout, "New branch created.\nSwitched to a new branch %s", args[1])
		}
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
	// checkoutCmd.Flags().StringP("branch", "b", "", "Create a new branch and switch to it")
	//checkoutCmd.LocalFlags().StringVarP(&branchName, "", "b", "", "branch name to create")

	// checkoutCmd.MarkPersistentFlagRequired("region")
}
