package webhook

import (
	"os"
	"os/signal"
	"sync"

	"net/http"

	config "github.com/sysincz/prometheus_webhook_logger/config"
	types "github.com/sysincz/prometheus_webhook_logger/types"

	logrus "github.com/Sirupsen/logrus"
)

var (
	log      = logrus.WithFields(logrus.Fields{"logger": "Webhook-server"})
	myConfig config.Config
)

func init() {
	// Set the log-level:
	logrus.SetLevel(logrus.DebugLevel)
}

//Run main function for webhook
func Run(myConfigFromMain config.Config, alertsChannel chan types.Alert, waitGroup *sync.WaitGroup) {
	log.Info("Webhook run")
	log.WithFields(logrus.Fields{"address": myConfigFromMain.WebhookAddress}).Info("Starting the Webhook server")

	// Populate the config:
	myConfig = myConfigFromMain

	// Set up a channel to handle shutdown:
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)

	// Listen for webhooks:
	http.ListenAndServe(myConfig.WebhookAddress, &Handler{AlertsChannel: alertsChannel})

	// Wait for shutdown:
	for {
		select {
		case <-signals:
			log.Info("Shutting down the Webhook server")

			// Tell main() that we're done:
			waitGroup.Done()
			return
		}
	}

}