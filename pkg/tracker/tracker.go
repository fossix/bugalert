// This package is used to get bugs and issues from different vendors
package tracker

import (
	"fmt"
)

type VendorType string

const (
	BUGZILLA VendorType = "bugzilla"
	NVBUG    VendorType = "NVIDIA"
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
	// case NVBUG:
	// 	return NewNVBug(conf)
	default:
		return nil, fmt.Errorf("Invalid vendor type")
	}

	return nil, nil
}

func SetRequestTimeout(timeout int) {
	globalTimeout = timeout
}
