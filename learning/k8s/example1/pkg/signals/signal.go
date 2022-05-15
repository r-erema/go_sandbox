package signals

import (
	"os"
	"os/signal"
	"syscall"
)

const osSignalChanSize = 2

func SetupSignalHandler() (stopCh <-chan struct{}) {
	var (
		onlyOneSignalHandler = make(chan struct{})
		shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
	)

	close(onlyOneSignalHandler)

	stop := make(chan struct{})
	signalCh := make(chan os.Signal, osSignalChanSize)

	signal.Notify(signalCh, shutdownSignals...)

	go func() {
		<-signalCh
		close(stop)
		<-signalCh
		os.Exit(1)
	}()

	return stop
}
