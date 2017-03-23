package helper

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kr/pretty"
)

var Log = logrus.New()

func InitLogging(isDebug bool, formatter logrus.Formatter) {
	Log.Formatter = formatter
	if isDebug {
		Log.Level = logrus.DebugLevel
	} else {
		Log.Level = logrus.InfoLevel
	}
}

// helper.FailOnError is on error behavior
func Fatal(err error, msg string) {
	if err != nil {
		Log.Fatalf("%s: %s, Error: %s", time.Now().Format(time.StampMilli), msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// LogError is log error behavior
func Error(err error, msg string) {
	if err != nil {
		Log.Errorf("%s: %s, Error: %s", time.Now().Format(time.StampMilli), msg, err)
	}
}

func Info(msg string) {
	Log.Infof("%s: %s", time.Now().Format(time.StampMilli), msg)
}

func Debug(msg string, obj interface{}) {
	Log.Debugf("%s: %s, %# v", time.Now().Format(time.StampMilli), msg, pretty.Formatter(obj))
}
