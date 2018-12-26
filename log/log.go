package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is a shared logger
var Logger = logrus.New()

func init() {
	Logger.Out = os.Stdout
}
