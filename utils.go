package itracker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// the default timeout in seconds
var globalTimeout int

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

func get(endpoint string, args map[string]string) ([]byte, error) {
	url := getURL(endpoint, args)

	timeout := time.Duration(globalTimeout) * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(http.StatusText(resp.StatusCode))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
