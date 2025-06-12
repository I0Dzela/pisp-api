package logger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/rifflock/lfshook"
	logrus "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"time"
)

const (
	typeFieldKey      string = "type"
	eventFieldKey     string = "event"
	reguestIdFieldKey string = "X-Request-Id"
)

const (
	rest   string = "rest"
	grpc   string = "grpc"
	socket string = "socket"
	pubsub string = "pubsub"
	db     string = "db"
)

var (
	GrpcEntry   = logrus.WithFields(logrus.Fields{typeFieldKey: grpc})
	RestEntry   = logrus.WithFields(logrus.Fields{typeFieldKey: rest})
	SocketEntry = logrus.WithFields(logrus.Fields{typeFieldKey: socket})
	PubsubEntry = logrus.WithFields(logrus.Fields{typeFieldKey: pubsub})
	DbEntry     = logrus.WithFields(logrus.Fields{typeFieldKey: db})
)

func RequestIdField(c context.Context) (key string, value interface{}) {
	requestId := "-"
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return reguestIdFieldKey, requestId
	}

	reqId := md.Get(reguestIdFieldKey)
	if len(reqId) > 0 {
		requestId = reqId[0]
	}

	return reguestIdFieldKey, requestId
}

func RequestIdFieldFromHeaders(c *gin.Context) (key string, value interface{}) {
	return reguestIdFieldKey, c.GetHeader(reguestIdFieldKey)
}

func EventField(e string) (key string, value interface{}) {
	return eventFieldKey, e
}

func WithRequestId(c *gin.Context) context.Context {
	return metadata.AppendToOutgoingContext(c.Request.Context(), reguestIdFieldKey, c.GetHeader(reguestIdFieldKey))
}

type LoggerX interface {
	WithFields(fields logrus.Fields) LoggerX
	Debug(m string)
	Info(m string)
	Warn(m string)
	Error(e *errors.Error)
	Fatal(e *errors.Error)
	Panic(e *errors.Error)
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(e *errors.Error, format string, args ...interface{})
	Fatalf(e *errors.Error, format string, args ...interface{})
	Panicf(e *errors.Error, format string, args ...interface{})
}

type loggerX struct {
	e *logrus.Entry
}

func (l loggerX) WithFields(fields logrus.Fields) LoggerX {
	return NewLogger(l.e.WithFields(fields))
}

func (l loggerX) Debug(m string) {
	l.e.Debug(m)
}

func (l loggerX) Info(m string) {
	l.e.Info(m)
}

func (l loggerX) Warn(m string) {
	l.e.Warn(m)
}

func (l loggerX) Error(e *errors.Error) {
	l.e.WithFields(logrus.Fields{"stack": string(e.Stack())}).Error(e.Error())
}

func (l loggerX) Fatal(e *errors.Error) {
	l.e.WithFields(logrus.Fields{"stack": string(e.Stack())}).Fatal(e.Error())
}

func (l loggerX) Panic(e *errors.Error) {
	l.e.WithFields(logrus.Fields{"stack": string(e.Stack())}).Panic(e.Error())
}

func (l loggerX) Debugf(format string, args ...interface{}) {
	l.e.Debugf(format, args...)
}

func (l loggerX) Infof(format string, args ...interface{}) {
	l.e.Infof(format, args...)
}

func (l loggerX) Warnf(format string, args ...interface{}) {
	l.e.Warnf(format, args...)
}

func (l loggerX) Errorf(e *errors.Error, format string, args ...interface{}) {
	l.e.WithFields(logrus.Fields{"stack": string(e.Stack())}).Errorf(format, args...)
}

func (l loggerX) Fatalf(e *errors.Error, format string, args ...interface{}) {
	l.e.WithFields(logrus.Fields{"stack": string(e.Stack())}).Fatalf(format, args...)
}

func (l loggerX) Panicf(e *errors.Error, format string, args ...interface{}) {
	l.e.WithFields(logrus.Fields{"stack": string(e.Stack())}).Panicf(format, args...)
}

func NewLogger(ex *logrus.Entry) LoggerX {
	return &loggerX{
		e: ex,
	}
}

func SetupLokiLogging(c *Config) error {
	hook, err := NewHook(c)
	if errors.Wrap(err, 0) != nil {
		return errors.Wrap(err, 0)
	}
	logrus.AddHook(hook)
	return nil
}

func SetupFileLogging(logFile string) {
	logrusFileHook := lfshook.NewHook(lfshook.PathMap{
		logrus.InfoLevel:  logFile,
		logrus.WarnLevel:  logFile,
		logrus.ErrorLevel: logFile,
		logrus.FatalLevel: logFile,
		logrus.PanicLevel: logFile,
	}, &logrus.TextFormatter{})
	logrus.AddHook(logrusFileHook)
}

func DefaultLogger(tag string) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}
		return fmt.Sprintf("[%s] %v |%s %3d %s| %13v | %15s | %s %-7s %s %s %-7s %s  %#v\n%s",
			tag,
			param.TimeStamp.Format("15:04:05 02.01.2006"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			"\033[90;47m", param.Request.Proto, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	})
}
