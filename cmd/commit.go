/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jplindgren/bastard-git/internal/object"
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
		repo := repository.GetRepository()
		if !repo.IsGitRepo() {
			fmt.Println("Not a git repository")
			os.Exit(1)
		}

		message := args[0]

		if message == "" {
			fmt.Println("Please provide a commit message")
			os.Exit(1)
		}

		fmt.Println("Creating commit with message: " + message)

		rootTree, err := repository.GenerateObjectTree(repo.WorkTree)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		//TODO: commit should check if a previous commit exists and set it as parent
		commit := object.NewCommit(repo.Email, message, rootTree.GetHash(), []string{})
		err = repo.Store.Write(commit)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		repo.UpdateIndex(rootTree)
		repo.UpdateRefHead(commit.ToString())
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
