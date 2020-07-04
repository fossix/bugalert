package cmd

import (
	"fmt"
	"strconv"
	"strings"

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

// Treats non-existent comment ID in the bug as empty comment and just returns
// the message.
func quoteComment(bug *itracker.Bug, commentID int, message string) string {
	var commentText string
	var commentCount int

	if err := bug.GetComments(); err != nil {
		panic(err)
	}

	for _, comment := range bug.Comments {
		if comment.ID == commentID {
			commentText = comment.Text
			commentCount = comment.Count
			break
		}
	}
	if commentText == "" {
		return message
	}

	newMessage := ""
	for _, line := range strings.Split(commentText, "\n") {
		newMessage = fmt.Sprintf("%s> %s\n", newMessage, line)
	}

	newMessage = fmt.Sprintf("(In reply to comment #%d)\n%s\n%s",
		commentCount, newMessage, message)

	return newMessage
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
	origmsg := message
	quote, _ := cmd.Flags().GetInt("quote")
	if quote != -1 {
		// Get a new message prefixed with the quote
		message = quoteComment(bug, quote, message)
	}

	editComment, _ := cmd.Flags().GetBool("edit")
	comment := []byte(message)
	if origmsg == "" || editComment {
		comment, err = editorInput(message, cfprefix)
		if err != nil {
			errLog(err)
		}
	}

	if len(comment) == 0 {
		fmt.Println("Empty comment message. Aborting.")
		return
	}

	public, _ := cmd.Flags().GetBool("public")

	var update itracker.BugUpdate
	update.Comment.Body = string(comment)
	update.Comment.Private = !public
	update.Comment.MarkDown = conf.doMarkdown

	if dryrun, _ := cmd.Flags().GetBool("dry-run"); dryrun == true {
		fmt.Println("Comment is public? ", !update.Comment.Private)
		fmt.Println("Markdown enabled? ", update.Comment.MarkDown)
		fmt.Println(update.Comment.Body)
	} else {
		bug.Update(&update)
	}
}
