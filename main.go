package main

import (
	"os"
	"strconv"

	"github.com/prodda/prodda/api"
	"github.com/prodda/prodda/registry"
	"github.com/prodda/prodda/schedule"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"gopkg.in/robfig/cron.v2"
)

var (
	username string
	password string
)

func main() {
	logger := lager.NewLogger("Prodda")
	sink := lager.NewReconfigurableSink(lager.NewWriterSink(os.Stdout, lager.DEBUG), lager.INFO)
	logger.RegisterSink(sink)

	portEnv := os.Getenv("PORT")
	port64, err := strconv.ParseUint(portEnv, 10, 0)
	port := uint(port64)
	if err != nil {
		logger.Fatal("Cannot parse port from environment", err, lager.Data{"PORT": portEnv})
	}

	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")

	logger.Info("Initializing registry")
	taskRegistry := registry.NewInMemoryTaskRegistry()
	logger.Info("Initializing registry complete")

	c := cron.New()
	handler := api.NewHandler(
		logger,
		username,
		password,
		taskRegistry,
		c)

	group := grouper.NewParallel(os.Kill, grouper.Members{
		grouper.Member{"schedule", schedule.NewRunner(c, logger)},
		grouper.Member{"api", api.NewRunner(port, handler, logger)},
	})
	process := ifrit.Invoke(group)

	logger.Info("Prodda started")
	err = <-process.Wait()
	if err != nil {
		logger.Fatal("Error running prodda", err)
	}
}
