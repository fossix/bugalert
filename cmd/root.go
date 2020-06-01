package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(versionCmd)

	listCmd.Flags().BoolP("all", "a", false, "Show all issues, don't restrict to user")
	listCmd.Flags().StringP("user", "u", "", "Show bugs associated with user")
	rootCmd.AddCommand(listCmd)

	showCmd.Flags().Bool("fuller", false, "Show more details of the bug shown")
	showCmd.Flags().Bool("fullest", false, "Show everything related to the bug")
	rootCmd.AddCommand(showCmd)

	rootCmd.AddCommand(historyCmd)
}

var rootCmd = &cobra.Command{
	Use:   "bugalert",
	Short: "Work with bugs and issues",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Bugalert v0.1")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
