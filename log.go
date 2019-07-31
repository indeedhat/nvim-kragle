package main

import (
	"fmt"
	"os"
)

var logPath string
var logFile *os.File

func initLog(path string) {
	logPath = path
	logOpen()
	log("New Client %s", config.ServerName)
}

func logOpen() {
	if "" == logPath {
		return
	}

	var err error
	logFile, err = os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if nil != err {
		return
	}
}

func logClose() {
	if nil == logFile {
		return
	}

	logFile.Close()
}

func log(message string, args ...interface{}) {
	logFile.WriteString(fmt.Sprintf(message, args...) + "\n")
}
