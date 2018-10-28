package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrAbortCommit    = errors.New("commit aborted")
	ErrContinueCommit = errors.New("continue with commit")
)

const warning = `
WARNING: ELF executable file(s) detected in commit

    You have several options -- you may:

      f) Attempt to fix the situation automatically (default)
      a) Abort and fix things manually
      c) Continue without change (ill advised)

So, what's it gonna be? [Fac] `

func chooseAction() error {
	fmt.Fprint(os.Stderr, warning)

	for {
		c, err := getch()
		if err != nil {
			return err
		}

		switch c {
		case 'F', 'f', '\r':
			return nil

		case 'A', 'a':
			fmt.Fprintln(os.Stderr, "\n\nAborting commit.\n")
			fmt.Fprintln(os.Stderr, "Please remove all binaries from the commit and update your .gitignore accordingly.")
			return ErrAbortCommit

		case 'C', 'c':
			fmt.Fprintln(os.Stderr, "\n\nContinuing binary commit.")
			fmt.Fprintln(os.Stderr, "NOTE: To bypass this prompt in the future, issue the following command:")
			fmt.Fprintln(os.Stderr, "\n   git config --bool precommit.allow-binaries 1")
			return ErrContinueCommit
		}
	}
}
