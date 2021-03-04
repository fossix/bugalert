package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

func editorInput(prefill, prefix string) ([]byte, error) {
	tempFile, err := ioutil.TempFile(os.TempDir(), prefix)
	if err != nil {
		return []byte{}, fmt.Errorf("Cannot create temporary file", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err = tempFile.Write([]byte(prefill)); err != nil {
		return []byte{}, err
	}

	if err = tempFile.Close(); err != nil {
		return []byte{}, err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		return []byte{}, fmt.Errorf("Please set your EDITOR environment variable to your favourite editor", err)
	}

	// This is for cases where $EDITOR might have arguments
	// like 'emacsclient -c' :)
	args := strings.Split(editor, " ")
	// Get the full executable path for the editor.
	args[0], err = exec.LookPath(args[0])
	if err != nil {
		return []byte{}, err
	}

	args = append(args, tempFile.Name())

	// We don't want to pass the 'editor' itself as a file argument
	editorExec, args := args[0], args[1:]

	cmd := exec.Command(editorExec, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		errLog(err)
	}

	content, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

// adapted from: https://gist.github.com/ik5/d8ecde700972d4378d87#gistcomment-3074524

var (
	Info    = Teal
	Warn    = Yellow
	Removed = Red
	Added   = Green
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
