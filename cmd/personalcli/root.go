package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "personalcli",
	Short: "A personal CLI to help with daily tasks",
	Long:  `A command-line tool written in Go to provide quick access to weather, todos, notes, and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no subcommand is given
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
