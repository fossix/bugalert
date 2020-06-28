package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/fossix/bugalert/pkg/itracker"
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Add comment bug/issue details",
	Long:  `Add/update comments and other details of bug/issue`,
	Args:  cobra.ExactArgs(1),
	Run:   addComment,
}

func addComment(cmd *cobra.Command, args []string) {
	conf := getConfig()
	bz := getBugzilla(conf)

	b, err := strconv.Atoi(args[0])
	if err != nil {
		panic(err)
	}
	bug, err := bz.GetBug(b)
	if err != nil {
		panic(err)
	}

	cfprefix := fmt.Sprintf("%d-comment*", b)
	if conf.doMarkdown {
		cfprefix = fmt.Sprintf("%s%s", cfprefix, ".md")
	}

	message, _ := cmd.Flags().GetString("message")

	// open editor for getting the comment
	// TODO this should be based on some condition, like 'quote this
	// comment', 'let me recheck', or an explicit 'open editor'
	comment, err := editorInput(message, cfprefix)
	if err != nil {
		errLog(err)
	}

	var update itracker.BugUpdate
	update.Comment.Body = string(comment)
	update.Comment.Private = true
	update.Comment.MarkDown = conf.doMarkdown

	bug.Update(&update)
}
