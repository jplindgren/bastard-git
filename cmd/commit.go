/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jplindgren/bastard-git/internal/repository"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record changes to the repository",
	Long: `A commit is a snapshot of the repository at that point in time. In original GIT a commit contains the content of the index (your stage).
	bGit is simpler, and does not contain the concept of the index. At the time of the commit we generate the tree of the objects commit, trees and blobs.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo := repository.GetRepository(cfg.user)
		if !repo.IsValid() {
			os.Exit(1)
		}

		message := args[0]

		if message == "" {
			fmt.Println("Please provide a commit message")
			os.Exit(1)
		}

		fmt.Println("Creating commit with message: " + message)

		// rootTree, err := repository.GenerateObjectTree(repo.WorkTree)
		// check(err)

		// //TODO: commit should check if a previous commit exists and set it as parent
		// branch, err := repo.GetCurrentBranch()
		// check(err)

		// commitHash, _, err := repo.GetBranchTip(branch)
		// check(err)

		// commit := object.NewCommit(repo.Email, message, rootTree.GetHash(), []string{commitHash})
		// err = repo.Store.Write(commit)
		// check(err)

		// repo.UpdateIndex(rootTree)
		// repo.UpdateRefHead(commit.ToString())
		if err := repo.Commit(message); err != nil {
			fmt.Fprintln(os.Stdout, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
