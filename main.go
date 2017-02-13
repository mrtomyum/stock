package main

import (
	//log "github.com/Sirupsen/logrus"
	//"github.com/sebest/logrusly"
	"github.com/mrtomyum/stock/ctrl"
)

//var logglyToken string = "3db1e177-a815-4887-a1f6-4a1a2b56b4b1"

func main() {
	//Log
	//logrus := log.New()
	//hook := logrusly.NewLogglyHook(logglyToken, "http://logs-01.loggly.com/inputs/", log.InfoLevel, "info")
	//logrus.Hooks.Add(hook)
	//defer hook.Flush()
	//log.WithFields(log.Fields{
	//	"name": "Tom NAVA Stock",
	//}).Info("Start Logrus")

	server := ctrl.Router()
	server.Run(":8001")
}