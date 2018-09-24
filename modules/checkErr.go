package storagemanagementcli

import "log"

var logFatal = log.Fatal

func checkErr(err error) {
	if err != nil {
		logFatal(err)
	}
}
