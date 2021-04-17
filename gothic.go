package gothic

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jrapoport/gothic/config"
	"github.com/jrapoport/gothic/core"
	"github.com/jrapoport/gothic/hosts"
	"github.com/jrapoport/gothic/utils"
)

// Main is the application main
func Main(c *config.Config) error {
	if c.IsDebug() {
		utils.PrettyPrint(c)
	}
	signalsToCatch := []os.Signal{
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGABRT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	}
	name := fmt.Sprintf("%s (%s)", c.Name, utils.ExecutableName())
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, signalsToCatch...)
	a, err := core.NewAPI(c)
	if err != nil {
		return err
	}
	defer func() {
		if err = a.Shutdown(); err != nil {
			c.Log().Error(err)
			return
		}
		c.Log().Infof("%s shut down", name)
	}()
	c.Log().Infof("starting %s...", name)
	err = hosts.Start(a, c)
	if err != nil {
		return err
	}
	defer func() {
		if err = hosts.Shutdown(); err != nil {
			c.Log().Error(err)
			return
		}
		c.Log().Infof("%s shut down", name)
	}()
	c.Log().Infof("%s %s started", name, c.Version())
	<-stopCh
	c.Log().Infof("%s shutting down", name)
	return nil
}
