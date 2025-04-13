package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/r-erema/go_sendbox/utils"
	"golang.org/x/sys/unix"
)

func main() {
	namespaceFDArg := os.Args[1]

	namespaceFD, err := strconv.Atoi(namespaceFDArg)
	if err != nil {
		log.Print(fmt.Errorf("string to ingeger conversion error: %w", err))
		os.Exit(1)
	}

	err = unix.Setns(namespaceFD, unix.CLONE_NEWNET)
	if err != nil {
		log.Print(
			fmt.Errorf("setting namespace for file descriptor `%d` error: %w", namespaceFD, err),
		)
		os.Exit(1)
	}

	inode, err := utils.NetworkNamespaceInodeNumber(unix.Getpid(), unix.Gettid())
	if err != nil {
		log.Print(fmt.Errorf("getting namespace inode number error: %w", err))
		os.Exit(1)
	}

	_, err = os.Stdout.WriteString(strconv.Itoa(inode))
	if err != nil {
		log.Print(fmt.Errorf("writing to STDOUT error: %w", err))
		os.Exit(1)
	}
}
