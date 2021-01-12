package cmd

import (
	"context"
	"time"

	"github.com/jrapoport/gothic/api"
	"github.com/jrapoport/gothic/conf"
	"github.com/jrapoport/gothic/rpc/servers"
	"github.com/jrapoport/gothic/storage"
	"github.com/jrapoport/gothic/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = cobra.Command{
	Use:  "serve",
	Long: "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, serve)
	},
}

func serve(globalConfig *conf.GlobalConfiguration, config *conf.Configuration) {
	ctx, err := api.WithConfig(context.Background(), config)
	if err != nil {
		logrus.Fatalf("Error loading instance config: %+v", err)
	}
	listenAndServe(ctx, globalConfig)
}

// listenAndServe starts the API servers
func listenAndServe(ctx context.Context, globalConfig *conf.GlobalConfiguration) {

	db, closeDB := openDB(globalConfig)
	defer closeDB()
	a := api.NewAPIWithVersion(ctx, globalConfig, db, Version)
	log := logrus.WithField("component", "api")

	api.ListenAndServeREST(a, globalConfig)
	servers.ListenAndServeRPC(a, globalConfig)

	done := make(chan struct{})
	defer close(done)
	util.WaitForTermination(log, done)

	log.Info("shutting down...")
}

func openDB(globalConfig *conf.GlobalConfiguration) (db *storage.Connection, closeDB func()) {
	// try a couple times to connect to the database
	var err error
	for i := 1; i <= 3; i++ {
		time.Sleep(time.Duration((i-1)*100) * time.Millisecond)
		db, err = storage.Dial(globalConfig)
		if err == nil {
			break
		}
		logrus.WithError(err).WithField("attempt", i).Warn("Error connecting to database")
	}
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}

	closeDB = func() {

	}

	return
}
