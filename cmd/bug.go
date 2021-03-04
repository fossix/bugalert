package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/fossix/bugalert/pkg/itracker"
)

type BugConfig struct {
	URL           string            `yaml:"url"`
	APIKey        string            `yaml:"api_key"`
	Users         []string          `yaml:"user_list"`
	DefaultUser   string            `yaml:"default_user"`
	TimeOut       int               `yaml:"timeout"`
	DefaultFilter string            `yaml:"default_filter"`
	Filters       map[string]string `yaml:"filters"`
	filtermap     map[string]string
	doMarkdown    bool
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

func bugDescription(bug *itracker.Bug) {
	fmt.Println()

	scanner := bufio.NewScanner(strings.NewReader(bug.Description.Text))
	for scanner.Scan() {
		fmt.Println("   ", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println()
	}
}

func print_indented(text string) {
	scanner := bufio.NewScanner(strings.NewReader(text))

	for scanner.Scan() {
		fmt.Println("   ", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func bugComments(bug *itracker.Bug) {
	for _, c := range bug.Comments {
		fmt.Printf("\n[#%d] On %s, %s wrote:\n",
			c.ID, c.CreationTime.Format("02/01/2006"), c.Creator)
		print_indented(c.Text)
	}
}

func getConfig() *BugConfig {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	f, err := ioutil.ReadFile(filepath.Join(u.HomeDir, ".bugalert.yml"))
	if err != nil {
		log.Fatal(err)
	}

	conf := BugConfig{}
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		log.Fatal(err)
	}

	return &conf
}

func getBugzilla(conf *BugConfig) *itracker.Tracker {
	timeout := 10

	bz, _ := itracker.NewTracker(itracker.BUGZILLA, conf.URL)
	bz.SetAPIKey(conf.APIKey)
	if conf.TimeOut != 0 {
		timeout = conf.TimeOut
	}
	bz.SetTimeout(timeout)

	return bz
}
