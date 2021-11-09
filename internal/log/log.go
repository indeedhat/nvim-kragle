package log

import (
	"fmt"
	"os"

	"github.com/indeedhat/nvim-kraggle/internal/config"
)

var logPath string
var logFile *os.File

// Init the loggger
func Init(path string) {
	logPath = path
	logOpen()

	Printf("New Client %s", config.Get().ServerName)
}

// Close the logger
func Close() {
	if logFile == nil {
		return
	}

	logFile.Close()
}

// Printf to the log file
//
// unlike Printf in the stdlib this will append a new line to the end of the string
func Printf(message string, args ...interface{}) {
	logFile.WriteString(fmt.Sprintf(message, args...) + "\n")
}

func logOpen() {
	if logPath == "" {
		return
	}

	logFile, _ = os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
}
