package main

import (
	"fmt"
	"os"
)

var logPath string
var logFile *os.File

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

func log(message string) {
	logFile.WriteString(fmt.Sprintln(message))
}
