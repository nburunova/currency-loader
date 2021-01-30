package log

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func tsNow() string {
	return time.Now().UTC().Format(time.StampMilli)
}

type Logger interface {
	logrus.FieldLogger
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"ts": tsNow(),
	})

	l.Logger.Infoln("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
		"ts":    tsNow(),
	})
}

type StructuredLogger struct {
	*logrus.Logger
}

func NewEmpty() *StructuredLogger {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return &StructuredLogger{logger}
}

func NewLogger(format, level string) *StructuredLogger {
	logger := logrus.New()
	logger.Out = os.Stdout
	f := strings.ToLower(format)
	switch f {
	case "json":
		logger.Formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		}
	case "text":
		logger.Formatter = &logrus.TextFormatter{
			ForceColors:     true,
			TimestampFormat: time.RFC3339Nano,
			FullTimestamp:   true,
		}
	default:
		logger.Warnf("log: invalid formatter: %v, continue with default", f)
	}

	l := strings.ToLower(level)
	sev, err := logrus.ParseLevel(l)
	if err != nil {
		logger.Warnf("log: invalid level: %v, continue with info", l)
		sev = logrus.InfoLevel
	}
	logger.Level = sev
	return &StructuredLogger{logger}
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	entry.Logger = entry.Logger.WithFields(logrus.Fields{
		"ts":         tsNow(),
		"user_agent": r.UserAgent(),
		"uri":        r.RequestURI,
	})
	entry.Logger.Infoln("request started")
	return entry
}
