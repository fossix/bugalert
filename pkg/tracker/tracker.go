// This package is used to get bugs and issues from different vendors
package tracker

import (
	"fmt"
)

type TrackerType string

const (
	BUGZILLA TrackerType = "bugzilla"
	GITHUB   TrackerType = "github"
)

type TrackerConfig struct {
	Url      string
	Endpoint string
	ApiKey   string
	Username string
	Password string
}

type Tracker interface {
	Get(string, map[string]string) ([]byte, error)
	Post(string, map[string]string, []byte) ([]byte, error)
	GetBug(id int) (*Bug, error)
	Search(map[string]string) ([]*Bug, error)
	GetUser(string) (*User, error)
}

func NewTracker(vendor TrackerType, conf TrackerConfig) (Tracker, error) {
	switch vendor {
	case BUGZILLA:
		return NewBugzilla(conf)
	default:
		return nil, fmt.Errorf("Invalid vendor type")
	}
}

func SetRequestTimeout(timeout int) {
	globalTimeout = timeout
}
