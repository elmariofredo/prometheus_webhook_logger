package logger

import (
	"os"
	"os/signal"
	"sync"

	config "github.com/sysincz/prometheus_webhook_logger/config"
	types "github.com/sysincz/prometheus_webhook_logger/types"

	logrus "github.com/Sirupsen/logrus"
)

var (
	log      = logrus.WithFields(logrus.Fields{"logger": "webhook-logger"})
	myConfig config.Config
)

func init() {
	// Set the log-level:
	logrus.SetLevel(logrus.DebugLevel)

}

//Run main function for send alert to logger
func Run(myConfigFromMain config.Config, alertsChannel chan types.Alert, waitGroup *sync.WaitGroup) {

	log.Info("Starting the Logger")

	// Populate the config:
	myConfig = myConfigFromMain

	// Set up a channel to handle shutdown:
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)

	// Handle incoming alerts:
	go func() {
		for {
			select {

			case alert := <-alertsChannel:

				log.WithFields(logrus.Fields{"status": alert.Status}).Debug("Received an alert")
				callLog(alert)
			}
		}
	}()

	// Wait for shutdown:
	for {
		select {
		case <-signals:
			log.Warn("Shutting down the logger")

			// Tell main() that we're done:
			waitGroup.Done()
			return
		}
	}

}
