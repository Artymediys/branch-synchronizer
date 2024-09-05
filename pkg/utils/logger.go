package utils

import (
	"log"
	"os"
)

func NewLogger() (*os.File, error) {
	logFile, err := os.OpenFile("./log/bsync.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	log.SetOutput(logFile)
	os.Stdout = logFile
	os.Stderr = logFile

	return logFile, err
}
