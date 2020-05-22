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
	URL    string `yaml:"url"`
	APIKey string `yaml:"api_key"`
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

	user, err := bz.GetUser("user@example.com")
	if err != nil {
		panic(err)
	}

	bugs, err := user.Bugs()
	if err != nil {
		panic(err)
	}

	for _, b := range bugs {
		if b.Status == "OPEN" || b.Status == "ASSIGNED" {
			fmt.Println(b.ID, b.Status, b.Summary)
		}
	}
}
