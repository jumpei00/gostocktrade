package log

import (
	"os"
	"github.com/sirupsen/logrus"
)


// SetLogging sets log using in this application
func SetLogging() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
}