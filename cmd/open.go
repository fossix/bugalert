package cmd

import (
	"fmt"
	"strconv"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open bug in browser",
	Args:  cobra.ExactArgs(1),
	Run:   openBug,
}

func openBug(cmd *cobra.Command, args []string) {
	b, err := strconv.Atoi(args[0])
	if err != nil {
		panic(err)
	}

	conf := getConfig()
	browser.OpenURL(fmt.Sprintf("%s/show_bug.cgi?id=%d", conf.URL, b))
}
