package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jrapoport/gothic/api"
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/rpc/servers"
	"github.com/jrapoport/gothic/storage"
	"github.com/segmentio/encoding/json"
	"github.com/spf13/cobra"
)

var configFile = ""

var cmdConfig *conf.Configuration

var rootCmd = &cobra.Command{
	Use:     "gothic",
	Version: conf.CurrentVersion(),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(cmd, args); err != nil {
			return err
		}
		return Main(cmdConfig)
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)
	rootCmd.AddCommand(serveCmd)
	pf := rootCmd.PersistentFlags()
	pf.StringVarP(&configFile, "config", "c", "", "the config file to use")
}

func initConfig(*cobra.Command, []string) error {
	c, err := conf.LoadConfiguration(configFile)
	if err != nil {
		return err
	}
	cmdConfig = c
	return nil
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

// Main is the application main
func Main(c *conf.Configuration) error {
	if conf.Debug {
		b, err := json.MarshalIndent(cmdConfig, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
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
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, signalsToCatch...)
	log := c.Log.WithField("logger", "gothic")
	db, err := storage.Dial(c, log)
	if err != nil {
		log.Fatal(err)
	}
	a := api.NewAPI(c, db)
	log.Info("starting gothic...")
	api.ListenAndServeREST(a, c)
	servers.ListenAndServeRPC(a, c)
	log.Info("gothic started")
	<-interruptChannel
	log.Info("gothic shutting down")
	return nil
}
