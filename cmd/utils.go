package cmd

import (
	"log"
	"strings"
)

func errLog(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// creates a map from a string of form "key:value"; multiple such pairs are
// separated by a semicolon
func makeFilter(fmap map[string]string, filterlist string) map[string]string {
	if filterlist == "" {
		return nil
	}

	if fmap == nil {
		fmap = make(map[string]string)
	}

	fpairs := strings.Split(filterlist, ";")
	for _, p := range fpairs {
		m := strings.Split(p, ":")
		if len(m) != 2 {
			continue
		}
		fmap[m[0]] = m[1]
	}
	return fmap
}
