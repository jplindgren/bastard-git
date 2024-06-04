package cmd

import (
	"github.com/jplindgren/bastard-git/internal/repository"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty Git repository",
	Run: func(cmd *cobra.Command, args []string) {
		repository.Init()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
