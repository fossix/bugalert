package cmd

import (
	"strings"
)

func errLog(err error) {
	if err != nil {
		panic(err)
	}
}

func errFatal(err error) {
	if err != nil {
		panic(err)
	}
}

func errWarn(err error) {
	if err != nil {
		panic(err)
	}
}

// creates a map from a string of form "key:value"; multiple such pairs are
// separated by a semicolon
func makeFilter(filterlist string) map[string]string {
	if filterlist == "" {
		return nil
	}

	fmap := make(map[string]string)
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
