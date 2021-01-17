package utils

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// WaitForTermination blocks until the system signals termination or done has a value
func WaitForTermination(log logrus.FieldLogger, done <-chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-signals:
		log.Infof("Triggering shutdown from signal %s", sig)
	case <-done:
		log.Infof("Shutting down...")
	}
}
