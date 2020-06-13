// This package is used to get bugs and issues from different vendors
package itracker

import (
	"fmt"
)

type VendorType string

const (
	BUGZILLA VendorType = "bugzilla"
)


type Tracker struct {
	url      string
	endpoint string
	apikey   string
	vendor   VendorType
	*Bugzilla
}

func NewTracker(vendor VendorType, url string) (*Tracker, error) {
	t := &Tracker{}

	switch vendor {
	case BUGZILLA:
		t.Bugzilla = NewBugzilla(url, "rest")
	default:
		return nil, fmt.Errorf("Invalid vendor type")
	}
	t.vendor = vendor

	return t, nil
}

func (track *Tracker) SetTimeout(timeout int) {
	globalTimeout = timeout
}
