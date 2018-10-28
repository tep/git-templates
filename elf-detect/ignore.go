package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

type gitIgnoreFile struct {
	name     string
	lines    []string
	mods     []string
	patterns map[string]bool
}

func loadGitIgnore(filename string) (*gitIgnoreFile, error) {
	gif := &gitIgnoreFile{name: filename, patterns: make(map[string]bool)}

	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return gif, nil
		}
		return nil, err
	}
	defer f.Close()

	scnr := bufio.NewScanner(f)

	for scnr.Scan() {
		line := scnr.Text()
		gif.lines = append(gif.lines, line)

		pat := []byte(line)
		for i := 0; i < len(pat); i++ {
			switch pat[i] {
			case '\\':
				i++
			case '#':
				pat = pat[:i]
			}
		}

		if pat = bytes.TrimSpace(pat); len(pat) > 0 {
			gif.patterns[string(pat)] = true
		}
	}

	if err := scnr.Err(); err != nil {
		return nil, err
	}

	return gif, nil
}

func (g *gitIgnoreFile) has(pat string) bool {
	return g.patterns[pat]
}

func (g *gitIgnoreFile) add(pat string) {
	g.mods = append(g.mods, pat)
	g.patterns[pat] = true
}

func (g *gitIgnoreFile) save() error {
	if len(g.mods) == 0 {
		return nil
	}

	fmt.Fprintln(os.Stderr, "  Updating .gitignore:")

	newname := g.name + ".new"
	f, err := os.Create(newname)
	if err != nil {
		return fmt.Errorf("creating file %q: %v", newname, err)
	}

	for _, l := range g.lines {
		fmt.Fprintln(f, l)
	}

	for _, m := range g.mods {
		fmt.Fprintf(os.Stderr, "    Adding: %s\n", m)
		fmt.Fprintln(f, m)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("closing file %q: %v", newname, err)
	}

	if err := os.Rename(newname, g.name); err != nil {
		return fmt.Errorf("renaming %q to %q: %v", newname, g.name, err)
	}

	fmt.Fprintf(os.Stderr, "  Staging file for commit: %s\n", g.name)

	if _, err := runGit("add", g.name); err != nil {
		return err
	}

	return nil
}
