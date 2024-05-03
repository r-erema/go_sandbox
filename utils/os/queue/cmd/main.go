package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"

	"github.com/r-erema/go_sendbox/utils/os/queue"
)

func main() {
	command := os.Args[1]

	switch command {
	case "send-to-queue":
		queueName := "/" + os.Args[2]
		sendToQueue(queueName)
	case "receive-from-queue":
		queueName := "/" + os.Args[2]
		receiveFromQueue(queueName)
	case "is-queue-existed":
		queueName := "/" + os.Args[2]
		isQueueExisted(queueName)
	case "close-queue":
		closeQueue()
	default:
		log.Fatal("unknown command")
	}
}

func sendToQueue(queueName string) {
	descriptor, err := queue.Open(queueName, syscall.O_RDWR|syscall.O_CREAT)
	if err != nil {
		log.Fatalf(fmt.Sprintf("open queue error %v", err))
	}

	data := os.Args[3]

	err = queue.Send(descriptor, []byte(data), 1)
	if err != nil {
		log.Fatalf(fmt.Sprintf("sending to queue erro %v", err))
	}

	_, err = os.Stdout.WriteString(fmt.Sprintf("OK(queue descr: %d)", descriptor))
	if err != nil {
		log.Fatalf(fmt.Sprintf("error writing to STDOUT: %v", err))
	}
}

func receiveFromQueue(queueName string) {
	descriptor, err := queue.Open(queueName, syscall.O_RDWR|syscall.O_CREAT)
	if err != nil {
		log.Fatalf(fmt.Sprintf("open queue error %v", err))
	}

	data, err := queue.Receive(descriptor)
	if err != nil {
		log.Fatalf(fmt.Sprintf("receiving from queue error %v", err))
	}

	_, err = os.Stdout.Write(append([]byte(fmt.Sprintf("Data from queue(descr: %d): ", descriptor)), data...))
	if err != nil {
		log.Fatalf(fmt.Sprintf("error writing to STDOUT: %v", err))
	}
}

func isQueueExisted(queueName string) {
	var msg string
	if queue.Exists(queueName) {
		msg = fmt.Sprintf("queue `%s` is existed", queueName)
	} else {
		msg = fmt.Sprintf("queue `%s` is not existed", queueName)
	}

	_, err := os.Stdout.WriteString(msg)
	if err != nil {
		log.Fatalf(fmt.Sprintf("error writing to STDOUT: %v", err))
	}
}

func closeQueue() {
	queueID := os.Args[2]

	id, err := strconv.Atoi(queueID)
	if err != nil {
		log.Fatalf(fmt.Sprintf("converting queue ID to int error %v", err))
	}

	err = queue.Close(id)
	if err != nil {
		log.Fatalf(fmt.Sprintf("closing queue error %v", err))
	}
}
