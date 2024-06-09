package cmd

import (
	"fmt"
	"os"

	"github.com/jplindgren/bastard-git/internal/repository"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Run: func(cmd *cobra.Command, args []string) {
		repo := repository.GetRepository(cfg.user)
		if !repo.IsValid() {
			os.Exit(1)
		}

		fmt.Println("Getting status")

		toBeCommited, err := repo.Diff()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(toBeCommited) == 0 {
			fmt.Println("Nothing to commit")
		} else {
			fmt.Println("Changes to be commited:")
			for _, file := range toBeCommited {
				fmt.Printf("%s:   %s\n", file.Status, file.Name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
