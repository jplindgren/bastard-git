/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type config struct {
	user string
}

var cfg config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bGit",
	Short: "This is a mini clone for Git made in GO",
	Long:  `This project is just for learning purposes and is not intended to be used in production`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	DisableFlagParsing: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfg.user, "user", os.Getenv("BGIT_USER"), "env variable BGIT_USER")
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
