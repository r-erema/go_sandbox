package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"syscall"
)

func main() {
	chrootPath := os.Args[1]

	response := struct {
		DirsBeforeChroot []string `json:"dirs_before_chroot"`
		DirsAfterChroot  []string `json:"dirs_after_chroot"`
	}{}

	files, err := filesInDir("/")
	if err != nil {
		log.Print(fmt.Errorf("reading files in dir before chroot error: %w", err))
		os.Exit(1)
	}

	response.DirsBeforeChroot = files

	if err = syscall.Chroot(chrootPath); err != nil {
		log.Print(fmt.Errorf("chroot error: %w", err))
		os.Exit(1)
	}

	files, err = filesInDir("/")
	if err != nil {
		log.Print(fmt.Errorf("reading files in dir after chroot error: %w", err))
		os.Exit(1)
	}

	response.DirsAfterChroot = files

	output, err := json.Marshal(response)
	if err != nil {
		log.Print(fmt.Errorf("marshaling error: %w", err))
		os.Exit(1)
	}

	_, err = os.Stdout.Write(output)
	if err != nil {
		log.Print(fmt.Errorf("writing to STDOUT error: %w", err))
		os.Exit(1)
	}
}

func filesInDir(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("reding dir error: %w", err)
	}

	files := make([]string, len(entries))

	for i := range entries {
		files[i] = entries[i].Name()
	}

	return files, nil
}
