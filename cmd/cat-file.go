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

// catFileCmd represents the cat-file command
var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "Provide content or type and size information for repository objects",
	Long: `Cat-file gets information about an object (commit, tree, or blob) by the hash.
The content is returned unzipped. Blobs contains the file content, trees contains the file names and hashes,
and commits contains the commit message and author information.`,

	Run: func(cmd *cobra.Command, args []string) {
		repo := repository.GetRepository(cfg.user)
		if !repo.IsValid() {
			os.Exit(1)
		}

		hash := args[0]

		content, err := repo.Store.Get(hash)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Stdout.Write([]byte(content))
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
