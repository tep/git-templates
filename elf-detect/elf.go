package main

import "os"

const elfMagic = string(0x7f) + "ELF"

func isElf(filename string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer f.Close()

	magic := make([]byte, 4)
	if _, err := f.Read(magic); err != nil {
		return false, err
	}

	return string(magic) == elfMagic, nil
}
