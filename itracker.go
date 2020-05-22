// This package is used to get bugs and issues from different vendors
package itracker

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type VendorType string

const (
	BUGZILLA VendorType = "bugzilla"
)

type Itracker interface{}

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

func makeargs(keys []string, values []string) (map[string]string, error) {
	if len(keys) != len(values) {
		return nil, fmt.Errorf("Length of keys and values don't match")
	}
	args := make(map[string]string)
	for i, k := range keys {
		args[k] = values[i]
	}

	return args, nil
}

func getURL(endpoint string, args map[string]string) string {
	log.Println(args)
	base, err := url.Parse(endpoint)
	if err != nil {
		return ""
	}

	params := url.Values{}
	for k, v := range args {
		params.Add(k, v)
	}

	base.RawQuery = params.Encode()

	return base.String()
}

// TODO: will have to handle http errors (404..)
func get(endpoint string, args map[string]string) ([]byte, error) {
	url := getURL(endpoint, args)

	log.Println(url)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
