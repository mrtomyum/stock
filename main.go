package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/sebest/logrusly"
	"github.com/mrtomyum/stock/controller"
)

var logglyToken string = "4cd7bdfb-0345-4205-aeee-53e85a030eda"

func main() {
	//Log
	logrus := log.New()
	hook := logrusly.NewLogglyHook(logglyToken, "http://logs-01.loggly.com/inputs/", log.InfoLevel, "info")
	logrus.Hooks.Add(hook)
	defer hook.Flush()
	log.WithFields(log.Fields{
		"name": "Tom NAVA Stock",
	}).Info("Start Logrus")

	server := controller.SetupRoute()
	server.Run(":8001")
}