package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func getch() (byte, error) {
	orig, err := unix.IoctlGetTermios(unix.Stdin, unix.TCGETS)
	if err != nil {
		return 0, err
	}

	ocp := *orig
	raw := &ocp

	raw.Iflag &^= uint32(unix.BRKINT | unix.IGNBRK | unix.ICRNL | unix.INLCR | unix.ISTRIP | unix.IXON | unix.IXOFF)
	raw.Lflag &^= uint32(unix.ECHO | unix.ECHONL | unix.ICANON | unix.IEXTEN)
	raw.Lflag |= uint32(unix.ISIG)
	raw.Cc[unix.VMIN] = 1
	raw.Cc[unix.VTIME] = 0

	if err := unix.IoctlSetTermios(unix.Stdin, unix.TCSETS, raw); err != nil {
		return 0, fmt.Errorf("setting raw input mode: %v", err)
	}

	in := make([]byte, 1)
	_, err = os.Stdin.Read(in)

	if err := unix.IoctlSetTermios(unix.Stdin, unix.TCSETS, orig); err != nil {
		return 0, fmt.Errorf("restoring normal input mode: %v", err)
	}

	return in[0], err
}
