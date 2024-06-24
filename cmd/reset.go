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

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset current HEAD to the specified state",
	Long:  `Set the current branch head (HEAD) to <commit>`,
	Run: func(cmd *cobra.Command, args []string) {
		repo := repository.GetRepository(cfg.user)
		if !repo.IsValid() {
			os.Exit(1)
		}

		commit := args[0]

		if commit == "" {
			fmt.Println("Please provide a commit hash")
			os.Exit(1)
		}

		err := repo.Reset(commit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not reset to commit: %s", commit)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
