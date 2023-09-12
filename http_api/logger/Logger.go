package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

const (
	colorDebug = 95
	colorInfo  = 94
	colorWarn  = 93
	colorError = 91
	colorPanic = 91
	colorFatal = 91
)

func NewLogger(name string, level int) (logger *Logger) {
	logger = &Logger{
		name:  name,
		level: level,
		log:   logrus.New(),
	}
	logger.log.SetOutput(os.Stdout)
	logger.log.SetLevel(logrus.TraceLevel)
	logger.log.SetReportCaller(true)
	logger.log.SetFormatter(&Formatter{
		HideKeys:        true,
		CallerFirst:     true,
		TimestampFormat: "2006-01-02 15:03:04",
		FieldsOrder:     []string{"message", "user_ip", "user_agent", "request_id"},
	})
	return
}

type Logger struct {
	log   *logrus.Logger
	name  string
	level int
}

func (l *Logger) handleCtx(args ...interface{}) (*logrus.Entry, []interface{}) {
	entry := l.log.WithFields(logrus.Fields{
		"request_id": "0",
		"user_ip":    "0.0.0.0",
		"user_agent": "unknown",
	})
	if len(args) > 0 {
		index := 0
		for i := 0; i < len(args); i++ {
			c, ok := args[i].(*gin.Context)
			if ok {
				index = i
				entry = l.log.WithFields(logrus.Fields{
					"request_id": c.GetHeader("request_id"),
					"user_ip":    c.GetHeader("user_ip"),
					"user_agent": c.GetHeader("User-Agent"),
				})
				break
			}
		}
		if index > 0 {
			args = append(args[:index], args[index+1:]...)
		} else {
			entry = l.log.WithFields(logrus.Fields{
				"request_id": "0",
				"user_ip":    "0.0.0.0",
				"user_agent": "unknown",
			})
		}
	}
	return entry, args
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	if l.level > LevelDebug {
		return
	}
	entry, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorDebug, l.name, fmt.Sprintf(format, args...))
	entry.Debug(msg)
}

func (l *Logger) Infof(format string, a ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	entry, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorInfo, l.name, fmt.Sprintf(format, args...))
	entry.Info(msg)
}

func (l *Logger) Warnf(format string, a ...interface{}) {
	if l.level > LevelWarn {
		return
	}
	entry, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorWarn, l.name, fmt.Sprintf(format, args...))
	entry.Warn(msg)
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	if l.level > LevelError {
		return
	}
	entry, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorError, l.name, fmt.Sprintf(format, args...))
	entry.Error(msg)
}

func (l *Logger) Panicf(format string, a ...interface{}) {
	if l.level > LevelPanic {
		return
	}
	entry, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorPanic, l.name, fmt.Sprintf(format, args...))
	entry.Panic(msg)
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	if l.level > LevelFatal {
		return
	}

	entry, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorFatal, l.name, fmt.Sprintf(format, args...))
	entry.Fatal(msg)
}

func (l *Logger) Debug(a ...interface{}) {
	if l.level > LevelDebug {
		return
	}
	entry, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorInfo, l.name, fmt.Sprintln(args...))
	entry.Debug(msg)
}

func (l *Logger) Info(a ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	entry, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorInfo, l.name, fmt.Sprintln(args...))
	entry.Info(msg)
}
func (l *Logger) Warn(a ...interface{}) {
	if l.level > LevelWarn {
		return
	}
	entry, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorWarn, l.name, fmt.Sprintln(args...))
	entry.Info(msg)
}

func (l *Logger) Error(a ...interface{}) {
	if l.level > LevelError {
		return
	}
	entry, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorError, l.name, fmt.Sprintln(args...))
	entry.Error(msg)
}

func (l *Logger) Panic(a ...interface{}) {
	if l.level > LevelPanic {
		return
	}
	entry, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorPanic, l.name, fmt.Sprintln(args...))
	entry.Panic(msg)
}

func (l *Logger) Fatal(a ...interface{}) {
	if l.level > LevelFatal {
		return
	}
	entry, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorFatal, l.name, fmt.Sprintln(args...))
	entry.Fatal(msg)
}
