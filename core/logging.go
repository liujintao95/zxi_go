package core

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logging = logrus.New()

func LogInit() {
	file, err := os.OpenFile("logs/log.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Logging.SetOutput(file)
		Logging.Out = file
	} else {
		Logging.Info("Failed to log to file_old, using default stderr")
	}

	Logging.Formatter = &logrus.TextFormatter{}
}
