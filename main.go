package main

import (
	"fmt"
	"github.com/fossix/itracker"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

type BugConfig struct {
	URL         string   `yaml:"url"`
	APIKey      string   `yaml:"api_key"`
	Users       []string `yaml:"user_list"`
	DefaultUser string   `yaml:"default_user"`
}

func main() {
	// Read our yaml config file
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	f, err := ioutil.ReadFile(filepath.Join(u.HomeDir, ".bugalert.yml"))
	if err != nil {
		panic(err)
	}

	var conf BugConfig
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		panic(err)
	}

	bz, _ := itracker.NewTracker(itracker.BUGZILLA, conf.URL)
	bz.SetAPIKey(conf.APIKey)
	bz.SetTimeout(10)

	user, err := bz.GetUser(conf.DefaultUser)
	if err != nil {
		panic(err)
	}

	bugs, err := user.Bugs()
	if err != nil {
		panic(err)
	}

	show_history := false
	for _, b := range bugs {
		if b.Status == "OPEN" || b.Status == "ASSIGNED" || b.Status == "NEEDINFO" {
			if b.Status == "NEEDINFO" && b.Flags[0].Requestee == user.Email {
				fmt.Println(b.Flags)
			}
			fmt.Println(b.ID, b.Summary)
			fmt.Printf(" %-25s %-25s\n",
				fmt.Sprintf("Status: %s", b.Status),
				fmt.Sprintf("Created On: %s", b.CreationTime.Format("02/01/2006")))
			if !show_history {
				continue
			}
			if err = b.GetHistory(); err != nil {
				fmt.Println(err)
				continue
			}

			var prev_when string
			var prev_who string
			for _, h := range b.History {
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
	}
}
