package main

import (
	"flag"
	"fmt"
	"os"
	sync "sync"

	config "github.com/sysincz/prometheus_webhook_logger/config"
	logger "github.com/sysincz/prometheus_webhook_logger/logger"
	types "github.com/sysincz/prometheus_webhook_logger/types"
	webhook "github.com/sysincz/prometheus_webhook_logger/webhook"

	logrus "github.com/Sirupsen/logrus"
)

var (
	//conf       config.Config
	log        = logrus.WithFields(logrus.Fields{"logger": "main"})
	waitGroup  = &sync.WaitGroup{}
	configFile = flag.String("config", "/config/logger.yaml", "The logger configuration file")
	debug      = flag.Bool("debug", false, "Set Log to debug level and print as text")
	//Usage show info
	Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
)

func init() {

	flag.Parse()

	if *debug {
		// The TextFormatter is default, you don't actually have to do this.
		logrus.SetFormatter(&logrus.TextFormatter{})
		//logrus.SetFormatter(&logrus.JSONFormatter{})
		// Set the log-level:
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		// Log as JSON instead of the default ASCII formatter.
		logrus.SetFormatter(&logrus.JSONFormatter{})
		// Set the log-level:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	log.Infof("Loading configuration file %s", *configFile)
	conf, _, err := config.LoadConfigFile(*configFile)
	if err != nil {
		log.Errorf("Error loading configuration: %s", err)
	}

	// Make sure we wait for everything to complete before bailing out:
	defer waitGroup.Wait()

	// Prepare a channel of events (to feed the digester):
	log.Info("Preparing the alerts channel")
	alertsChannel := make(chan types.Alert)

	// Prepare to have background GoRoutines running:
	waitGroup.Add(1)

	// Start webhook server:
	go webhook.Run(*conf, alertsChannel, waitGroup)

	// Start the logger:
	go logger.Run(*conf, alertsChannel, waitGroup)

}
