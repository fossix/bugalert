package cmd

import (
	"fmt"
	"github.com/fossix/itracker"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strconv"

	"strings"
)

type BugConfig struct {
	URL         string   `yaml:"url"`
	APIKey      string   `yaml:"api_key"`
	Users       []string `yaml:"user_list"`
	DefaultUser string   `yaml:"default_user"`
	TimeOut     int      `yaml:"timeout"`
}

func bugSummary(bug *itracker.Bug) {
	fmt.Println(bug.ID, bug.Summary)
	fmt.Println("Status:", bug.Status)
	fmt.Println("Priority:", bug.Priority)
	fmt.Println("Severity:", bug.Severity)
	fmt.Println("Created on:", bug.CreationTime.Format("02/01/2006"))
	fmt.Printf("Creator: %s <%s>\n",
		bug.Creator.RealName, bug.Creator.Email)

	fmt.Printf("Assigned to: %s <%s>\n",
		bug.AssignedTo.RealName, bug.AssignedTo.Email)

	fmt.Printf("QA Contact: %s <%s>\n",
		bug.QaContact.RealName, bug.QaContact.Email)
}

func getConfig() (*BugConfig, error) {
	u, err := user.Current()
	if err != nil {
		return nil, (err)
	}

	f, err := ioutil.ReadFile(filepath.Join(u.HomeDir, ".bugalert.yml"))
	if err != nil {
		return nil, err
	}

	conf := BugConfig{}
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		panic(err)
	}

	return &conf, nil
}

func getBugzilla(conf *BugConfig) (*itracker.Tracker, error) {
	timeout := 10

	bz, _ := itracker.NewTracker(itracker.BUGZILLA, conf.URL)
	bz.SetAPIKey(conf.APIKey)
	if conf.TimeOut != 0 {
		timeout = conf.TimeOut
	}
	bz.SetTimeout(timeout)

	return bz, nil
}

func listBug(cmd *cobra.Command, args []string) {
	conf, _ := getConfig()
	bz, _ := getBugzilla(conf)

	username, err := cmd.Flags().GetString("user")
	if err != nil {
		panic(err)
	}
	if username == "" {
		username = conf.DefaultUser
	}
	if username == "" {
		fmt.Println("Warning: default_user config or --user option is not provided. Fetching all items")
	}

	user, err := bz.GetUser(username)
	if err != nil {
		panic(err)
	}

	bugs, err := user.Bugs()
	if err != nil {
		panic(err)
	}

	for _, b := range bugs {
		fmt.Println(b.ID, b.Summary)
	}
}

func showBug(cmd *cobra.Command, args []string) {
	conf, _ := getConfig()
	bz, _ := getBugzilla(conf)

	b, err := strconv.Atoi(args[0])
	if err != nil {
		panic(err)
	}
	bug, err := bz.GetBug(b)
	if err != nil {
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
	if !fuller {
		return
	}

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

func showHistory(cmd *cobra.Command, args []string) {
	conf, _ := getConfig()
	bz, _ := getBugzilla(conf)

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
	var prev_when string
	var prev_who string
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
			if c.Added != "" {
				fmt.Println("      + ", c.Added)
			}
			if c.Removed != "" {
				fmt.Println("      - ", c.Removed)
			}
		}
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all bugs/issues",
	Long: `This command lists all the bugs and issues. If default_user
config option is set to a user id, then only bugs associated with
that user is listed. This can be overridden with --user option.`,
	Run: listBug,
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show bug/issue details",
	Long: `This command shows the given bug's details. More details can
be obtained with --fuller/--fullest options.`,
	Args: cobra.ExactArgs(1),
	Run:  showBug,
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "show bug history",
	Long:  `This command shows the given bug's history.`,
	Args:  cobra.ExactArgs(1),
	Run:   showHistory,
}
