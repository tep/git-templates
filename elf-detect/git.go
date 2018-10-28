package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ErrBinariesAllowed = errors.New("binaries allowed")

func repoRoot() (string, error) {
	v, err := runGit("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}

	return filepath.Clean(strings.TrimSpace(string(v))), nil
}

func allowBinaries() error {
	v, err := runGit("config", "--bool", "precommit.allow-binaries")
	if err != nil && err.msg != "" {
		return err
	}

	if string(v) == "true" {
		return ErrBinariesAllowed
	}

	return nil
}

func binsAdded() ([]string, error) {
	added, err := addedFiles()
	if err != nil {
		return nil, err
	}

	if len(added) == 0 {
		return nil, nil
	}

	var bins []string
	for _, a := range added {
		ok, err := isElf(a)
		if err != nil {
			return nil, err
		}

		if ok {
			bins = append(bins, a)
		}
	}

	return bins, nil
}

func addedFiles() ([]string, error) {
	v, err := runGit("diff", "--cached", "--name-only", "--diff-filter=A", "-z")
	if err != nil {
		return nil, err
	}

	var list []string

	for _, a := range bytes.Split(v, []byte{0}) {
		if len(a) > 0 {
			list = append(list, string(a))
		}
	}

	return list, nil
}

func unstage(filename string) error {
	fmt.Fprintf(os.Stderr, "  Unstaging file %q from commit... ", filename)

	_, err := runGit("rm", "--quiet", "--cached", filename)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "OK")
	return nil
}

func runGit(args ...string) ([]byte, *gitError) {
	cmd := exec.Command("git", args...)
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err == nil {
		return out, nil
	}

	var msg string
	if xerr, ok := err.(*exec.ExitError); ok {
		msg = string(xerr.Stderr)
	}

	return nil, &gitError{strings.Join(cmd.Args, " "), msg, err}
}

type gitError struct {
	cmd string
	msg string
	err error
}

func (e *gitError) Error() string {
	m := e.err.Error()
	if e.msg != "" {
		m = e.msg
	}

	return fmt.Sprintf("command=%q: %s", e.cmd, m)
}
