// This package is used to get bugs and issues from different vendors
package tracker

import (
	"fmt"
)

type VendorType string

const (
	BUGZILLA VendorType = "bugzilla"
	GITHUB   VendorType = "Github"
)

type TrackerConfig struct {
	Url      string
	Endpoint string
	ApiKey   string
}

type Tracker interface {
	GetBug(id int) (*Bug, error)
	Search(map[string]string) ([]*Bug, error)
	GetUser(string) (*User, error)
}

func NewTracker(vendor VendorType, conf TrackerConfig) (Tracker, error) {
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
