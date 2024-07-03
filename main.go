package main

import (
	"fmt"
	"os"

	"github.com/barretodotcom/golang-git/command"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git",
	Short: "gogit is a minimalistic git cli tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func main() {
	rootCmd.AddCommand(command.Init)
	rootCmd.AddCommand(command.Add)
	rootCmd.AddCommand(command.Status)
	rootCmd.AddCommand(command.Commit)
	rootCmd.AddCommand(command.RollbackFile)
	rootCmd.AddCommand(command.Reset)
	rootCmd.AddCommand(command.Rollback)
	rootCmd.AddCommand(command.Log)
	rootCmd.AddCommand(command.Show)

	err := rootCmd.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
