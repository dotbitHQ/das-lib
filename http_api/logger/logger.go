package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"runtime/debug"
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
	colourDebug = 0
	colourInfo  = 94
	colourWarn  = 93
	colourError = 91
	colourPanic = 91
	colourFatal = 91
)

type Logger struct {
	log   *zap.SugaredLogger
	name  string
	level int
}

func (l *Logger) ErrStack() {
	l.Warn(string(debug.Stack()))
}

type ReqeustInfo struct {
	RequestId string
	UserIp    string
	UserAgent string
}

func (l *Logger) handleCtx(args ...interface{}) (ReqeustInfo, []interface{}) {
	requestInfo := ReqeustInfo{
		RequestId: "0",
		UserIp:    "0.0.0.0",
		UserAgent: "unknown",
	}
	if len(args) > 0 {
		index := 0
		for i := 0; i < len(args); i++ {
			c, ok := args[i].(context.Context)
			if ok {
				index = i
				if requestID, ok := c.Value("request_id").(string); ok {
					requestInfo.RequestId = requestID
				}
				if userAgent, ok := c.Value("user_agent").(string); ok {
					requestInfo.UserAgent = userAgent
				}
				if userIp, ok := c.Value("user_ip").(string); ok {
					requestInfo.UserIp = userIp
				}
				break
			}
		}
		if index > 0 {
			args = append(args[:index], args[index+1:]...)
		}
	}
	return requestInfo, args
}
func (l *Logger) Debugf(format string, a ...interface{}) {
	if l.level > LevelDebug {
		return
	}
	res, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm [DEBUG]  [%s] [%s] [%s] ▶ [%s] %s \x1b[0m \n", colourDebug, res.RequestId, res.UserIp, res.UserAgent, l.name, fmt.Sprintf(format, args...))
	l.log.Debug(msg)
}

func (l *Logger) Infof(format string, a ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	res, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm [INFO]  [%s] [%s] [%s] ▶ [%s] %s \x1b[0m \n", colourInfo, res.RequestId, res.UserIp, res.UserAgent, l.name, fmt.Sprintf(format, args...))
	l.log.Info(msg)
}

func (l *Logger) Warnf(format string, a ...interface{}) {
	if l.level > LevelWarn {
		return
	}
	res, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm [WARN]  [%s] [%s] [%s] ▶ [%s] %s \x1b[0m \n", colourWarn, res.RequestId, res.UserIp, res.UserAgent, l.name, fmt.Sprintf(format, args...))
	l.log.Warn(msg)
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	if l.level > LevelError {
		return
	}
	res, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm [ERROR]  [%s] [%s] [%s] ▶ [%s] %s \x1b[0m \n", colourError, res.RequestId, res.UserIp, res.UserAgent, l.name, fmt.Sprintf(format, args...))
	l.log.Error(msg)
}

func (l *Logger) Panicf(format string, a ...interface{}) {
	if l.level > LevelPanic {
		return
	}
	res, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm [PANIC]  [%s] [%s] [%s] ▶ [%s] %s \x1b[0m \n", colourPanic, res.RequestId, res.UserIp, res.UserAgent, l.name, fmt.Sprintf(format, args...))
	l.log.Panic(msg)
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	if l.level > LevelFatal {
		return
	}
	res, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm [FATAl]  [%s] [%s] [%s] ▶ [%s] %s \x1b[0m \n", colourFatal, res.RequestId, res.UserIp, res.UserAgent, l.name, fmt.Sprintf(format, args...))
	l.log.Fatal(msg)
}

func (l *Logger) Debug(a ...interface{}) {
	if l.level > LevelDebug {
		return
	}
	res, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm [DEBUG]  [%s] [%s] [%s] ▶ [%s] %s \x1b[0m \n", colourDebug, res.RequestId, res.UserIp, res.UserAgent, l.name, fmt.Sprintln(args...))

	l.log.Debug(msg)
}

func (l *Logger) Info(a ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	res, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm [INFO]  [%s] [%s] [%s] ▶ [%s] %s \x1b[0m \n", colourInfo, res.RequestId, res.UserIp, res.UserAgent, l.name, fmt.Sprintln(args...))
	l.log.Info(msg)
}

func (l *Logger) Warn(a ...interface{}) {
	if l.level > LevelWarn {
		return
	}
	res, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm [WARN]  [%s] [%s] [%s] ▶ [%s] %s \x1b[0m \n", colourWarn, res.RequestId, res.UserIp, res.UserAgent, l.name, fmt.Sprintln(args...))
	l.log.Warn(msg)
}

func (l *Logger) Error(a ...interface{}) {
	if l.level > LevelError {
		return
	}
	res, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm [ERROR]  [%s] [%s] [%s] ▶ [%s] %s \x1b[0m \n", colourError, res.RequestId, res.UserIp, res.UserAgent, l.name, fmt.Sprintln(args...))
	l.log.Error(msg)
}

func (l *Logger) Panic(a ...interface{}) {
	if l.level > LevelPanic {
		return
	}
	res, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm [PANIC]  [%s] [%s] [%s] ▶ [%s] %s \x1b[0m \n", colourPanic, res.RequestId, res.UserIp, res.UserAgent, l.name, fmt.Sprintln(args...))
	l.log.Panic(msg)
}

func (l *Logger) Fatal(a ...interface{}) {
	if l.level > LevelFatal {
		return
	}
	res, args := l.handleCtx(a...)
	msg := fmt.Sprintf("\x1b[%dm [FATAL]    [%s] [%s] [%s] ▶ [%s] %s \x1b[0m \n", colourFatal, res.RequestId, res.UserIp, res.UserAgent, l.name, fmt.Sprintln(args...))
	l.log.Fatal(msg)
}
