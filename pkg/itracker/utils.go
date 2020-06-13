package itracker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

// create a URL by adding the maps as request params. If the values in the map
// are separated by a "|", then a separate param will be added for the same key.
func getURL(endpoint string, args map[string]string) (string, error) {
	base, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	for key, value := range args {
		values := strings.Split(value, "|")
		for _, v := range values {
			params.Add(key, v)
		}
	}

	base.RawQuery = params.Encode()

	return base.String(), nil
}

func get(endpoint string, args map[string]string) ([]byte, error) {
	url, err := getURL(endpoint, args)
	if err != nil {
		return nil, err
	}

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
