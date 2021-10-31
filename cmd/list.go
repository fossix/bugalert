package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/fossix/bugalert/pkg/tracker"
)

func listBug(cmd *cobra.Command, args []string) {
	var bugs []*tracker.Bug
	var err error

	conf := getConfig()
	bz := getTracker(conf)
	username := conf.DefaultUser
	filter := conf.DefaultFilter

	if _filter, err := cmd.Flags().GetString("filter"); err == nil {
		if _filter != "" {
			filter = _filter
		}
	}

	if filter_name, err := cmd.Flags().GetString("by-filter"); err == nil {
		ok := true
		if filter_name != "" {
			if filter, ok = conf.Filters[filter_name]; !ok {
				fmt.Println("Predefined filter", filter_name,
					"not found, defined filters are:\n")
				for k, v := range conf.Filters {
					fmt.Printf("    %s: %s\n", k, v)
				}
				return
			}
		}
	}

	if skipfilter, _ := cmd.Flags().GetBool("nofilter"); skipfilter == false {
		conf.filtermap = makeFilter(conf.filtermap, filter)
	}

	order, _ := cmd.Flags().GetBool("order")
	// I don't see the bugs really sorted with the 'order' parameter. Should
	// be fairly easy like the code below if it does.
	//
	// order := fmt.Sprintf("order:%s", order_field)
	// conf.filtermap = makeFilter(conf.filtermap, order)

	if _username, err := cmd.Flags().GetString("user"); err == nil {
		if _username != "" {
			username = _username
		}
	}

	allusers, _ := cmd.Flags().GetBool("all")
	if username != "" && allusers == false {
		user, err := bz.GetUser(username)
		errLog(err)

		bugs, err = user.Bugs(conf.filtermap)
		errLog(err)
	} else {
		// user didn't specify '-a' explicitly, so print a warning
		if allusers == false {
			fmt.Println("Warning: default_user config or --user option not provided. Fetching all items")
		}
		bugs, err = bz.Search(conf.filtermap)
		errLog(err)
	}

	if limit, _ := cmd.Flags().GetInt("limit"); limit != 0 {
		if limit > len(bugs) {
			limit = len(bugs)
		}
		bugs = bugs[len(bugs)-limit : len(bugs)]
	}

	if order == true {
		sort.Slice(bugs, func(i, j int) bool {
			return bugs[i].LastChangeTime.Before(bugs[j].LastChangeTime)
		})
	}

	for _, b := range bugs {
		fmt.Printf("%7d [ %-10s] %s\n", b.ID, b.Status, b.Summary)
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all bugs/issues",
	Long: `Lists all the bugs and issues. If default_user
config option is set to a user id, then only bugs associated with
that user is listed. This can be overridden with --user option.`,
	Run: listBug,
}
