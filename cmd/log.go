package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "show bug history",
	Long:  "Display given bug's history.",
	Args:  cobra.ExactArgs(1),
	Run:   showLog,
}

func showLog(cmd *cobra.Command, args []string) {
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

	err = bug.GetHistory()
	if err != nil {
		panic(err)
	}

	bugSummary(bug)
	var (
		prev_when,
		prev_who string
	)
	for _, h := range bug.History {
		when := h.When.Format("2006, January 02")
		who := h.Who
		if when != prev_when && who != prev_who {
			fmt.Println("  On", when, "by", who)
		}
		prev_when = when
		prev_who = who

		for _, c := range h.Changes {
			fmt.Printf("    %s: \n", c.FieldName)
			if c.Removed != "" {
				fmt.Println(Removed("      - ", c.Removed))
			}
			if c.Added != "" {
				fmt.Println(Added("      + ", c.Added))
			}

		}
	}
}
