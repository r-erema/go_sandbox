package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/r-erema/go_sendbox/utils"
	"k8s.io/apimachinery/pkg/util/json"
)

func main() {
	sourceDir, targetDir, sourceDirFileName, sourceDirFileContent, outputFD := os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5]

	fd, err := strconv.Atoi(outputFD)
	if err != nil {
		log.Panicf("string to int conversion error: %s", err)
	}

	outputSocket := os.NewFile(uintptr(fd), "")
	defer func() {
		err = outputSocket.Close()
		if err != nil {
			log.Panicf("closing output socket error: %s", err)
		}
	}()

	err = os.MkdirAll(sourceDir, 0o755)
	if err != nil {
		log.Panicf("creation dir `%s` error: %s", sourceDir, err)
	}

	err = os.MkdirAll(targetDir, 0o755)
	if err != nil {
		log.Panicf("creation dir `%s` error: %s", targetDir, err)
	}

	fileInSource := filepath.Join(sourceDir, sourceDirFileName)
	f, err := os.Create(fileInSource)
	if err != nil {
		log.Panicf("creation file `%s` error: %s", sourceDirFileName, err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("closing file error: %s", err)
		}
	}()

	_, err = f.WriteString(sourceDirFileContent)
	if err != nil {
		log.Panicf("writing string `%s` into file `%s` error: %s", sourceDirFileContent, sourceDirFileName, err)
	}

	err = syscall.Mount(sourceDir, targetDir, "", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		log.Panicf("mounting `%s` into `%s` error: %s", sourceDir, targetDir, err)
	}

	fileInTarget := filepath.Join(targetDir, sourceDirFileName)
	data, err := os.ReadFile(fileInTarget)
	if err != nil {
		log.Panicf("reading file `%s` error: %s", fileInTarget, err)
	}

	files, err := utils.FilesInDir(targetDir)
	if err != nil {
		log.Panicf("getting files list in directory `%s` error: %s", targetDir, err)
	}

	ns, err := utils.MountNamespaceInodeNumber(os.Getpid())
	if err != nil {
		log.Panicf("getting mount namespace error: %s", err)
	}

	output := struct {
		PID                int      `json:"pid"`
		MountNamespace     int      `json:"mount_namespace"`
		File               string   `json:"file"`
		FileContent        string   `json:"fileContent"`
		TargetDirFilesList []string `json:"targetDirFilesList"`
	}{
		PID:                os.Getpid(),
		MountNamespace:     ns,
		File:               fileInTarget,
		FileContent:        string(data),
		TargetDirFilesList: files,
	}

	data, err = json.Marshal(output)
	if err != nil {
		log.Panicf("json marshaling error: %s", err)
	}

	_, err = outputSocket.Write(data)
	if err != nil {
		log.Panicf("writing to STDOUT error: %s", err)
	}
}
