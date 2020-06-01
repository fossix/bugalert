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
)

type BugConfig struct {
	URL         string   `yaml:"url"`
	APIKey      string   `yaml:"api_key"`
	Users       []string `yaml:"user_list"`
	DefaultUser string   `yaml:"default_user"`
	TimeOut     int      `yaml:"timeout"`
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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all bugs/issues",
	Long: `This command lists all the bugs and issues. If default_user
config option is set to a user id, then only bugs associated with
that user is listed. This can be overridden with --user option.`,
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show bug/issue details",
	Long: `This command shows the given bug's details. More details can be
	obtained with --fuller option.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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

		fmt.Println(bug.Summary)

	},
}
