package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show bug/issue details",
	Long: `Show the given bug's details. More details can
be obtained with --fuller/--fullest options.`,
	Args: cobra.ExactArgs(1),
	Run:  showBug,
}

func showBug(cmd *cobra.Command, args []string) {
	conf := getConfig()
	bz := getTracker(conf)

	b, err := strconv.Atoi(args[0])
	if err != nil {
		panic(err)
	}
	bug, err := bz.GetBug(b)
	if err != nil {
		panic(err)
	}

	if err = bug.GetComments(); err != nil {
		panic(err)
	}

	fullest, _ := cmd.Flags().GetBool("fullest")
	if fullest {
		// Should do a pretty json dump
		fmt.Printf("%+v\n", bug)
		return
	}

	bugSummary(bug)
	fuller, _ := cmd.Flags().GetBool("fuller")
	if fuller {
		fmt.Println("Resolution:", bug.Resolution)
		fmt.Print("Blocks: ")
		if len(bug.Blocks) > 0 {
			fmt.Println(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(bug.Blocks)), ", "), "[]"))
		} else {
			fmt.Println("None")
		}

		fmt.Print("Depends on: ")
		if len(bug.DependsOn) > 0 {
			fmt.Println(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(bug.DependsOn)), ", "), "[]"))
		} else {
			fmt.Println("None")
		}

		if bug.DupeOf != 0 {
			fmt.Println("Duplicate of:", bug.DupeOf)
		}

		if len(bug.Cc) > 0 {
			fmt.Println("CC:", strings.Join(bug.Cc, ", "))
		}
	}

	bugDescription(bug)
	if comments, _ := cmd.Flags().GetBool("comments"); comments {
		bugComments(bug)
	}
}
