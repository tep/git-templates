// gpch - git pre-commit hook
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	err := run()

	switch err {
	case nil, ErrBinariesAllowed, ErrContinueCommit:
		os.Exit(0)

	case ErrAbortCommit:
		os.Exit(1)

	default:
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	root, err := repoRoot()
	if err != nil {
		return err
	}

	if err := allowBinaries(); err != nil {
		return err
	}

	bins, err := binsAdded()
	if err != nil || len(bins) == 0 {
		return err
	}

	if err := chooseAction(); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "\n\nAttempting to auto-fix:")

	gif, err := loadGitIgnore(filepath.Join(root, ".gitignore"))
	if err != nil {
		return err
	}

	for _, bin := range bins {
		if err := unstage(bin); err != nil {
			return err
		}

		if bb := filepath.Base(bin); !gif.has(bb) {
			gif.add(bb)
		}
	}

	if err := gif.save(); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "auto-fix completed")
	return nil
}
