package logger

import (
	cl "github.com/I0Dzela/pisp-common/logger"
	"github.com/I0Dzela/pisp-specs/config"
	"github.com/sirupsen/logrus"
)

func Init() error {
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true, FullTimestamp: true})
	logrus.SetLevel(logrus.Level(config.Server.LogLevel))

	return nil
}

func NewLogger(ex *logrus.Entry) cl.LoggerX {
	e := logrus.WithFields(logrus.Fields{"service": "specs"})
	if ex != nil {
		for k := range ex.Data {
			e.Data[k] = ex.Data[k]
		}
	}

	return cl.NewLogger(e)
}
