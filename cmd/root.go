package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)

	// the 'list' command
	listCmd.Flags().BoolP("all", "a", false, "Show all issues, don't restrict to user")
	listCmd.Flags().StringP("user", "u", "", "Show bugs associated with user")
	listCmd.Flags().StringP("filter", "", "", "Show bugs filtered by the given condition")
	listCmd.Flags().BoolP("nofilter", "", false, "Don't filter bugs, skip default filter too")
	listCmd.Flags().Int("limit", 0, "List bugs sort by ascending, limiting to N")
	// Bugzilla doesn't seem to send a sorted array. We will do it manually
	// and only sort by last_change_time. If bugzilla works, then the choice
	// can be given to the user on the fields to sort.
	// listCmd.Flags().String("order", "last_change_time", "Sort the bugs by the given field [Default is last_change_time]")
	listCmd.Flags().Bool("order", false, "Sort the bugs by last changed time")
	listCmd.Flags().StringP("by-filter", "", "", "Show bugs filtered by a predefined filter name")
	rootCmd.AddCommand(listCmd)

	// The 'show' command
	showCmd.Flags().Bool("fuller", false, "Show more details of the bug shown")
	showCmd.Flags().Bool("fullest", false, "Show everything related to the bug")
	showCmd.Flags().Bool("comments", false, "Show comments on this bug")
	rootCmd.AddCommand(showCmd)

	// The 'comment' command
	commentCmd.Flags().Bool("edit", false, "Open EDITOR to edit comment")
	commentCmd.Flags().Bool("public", false, "Posted comment will be public")
	commentCmd.Flags().Bool("dry-run", false, "Don't do the actual update")
	commentCmd.Flags().IntP("quote", "q", -1, "Quote the provided comment number")
	commentCmd.Flags().StringP("message", "m", "", "Add comment to bug")
	rootCmd.AddCommand(commentCmd)

	// The 'log' Command
	rootCmd.AddCommand(logCmd)

	// Open bug in browser
	rootCmd.AddCommand(openCmd)
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
