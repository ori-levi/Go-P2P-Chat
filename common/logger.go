package common

import (
	"os"

	"github.com/withmandala/go-log"
)

var Logger = NewLogger()

func NewLogger() *log.Logger {
	return log.New(os.Stderr).WithColor().WithDebug().WithTimestamp()
}
