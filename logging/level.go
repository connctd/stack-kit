package logging

import (
	"github.com/go-kit/kit/log"
)

// Debug returns a logger that includes a severity of level debug
func Debug(logger log.Logger) log.Logger {
	return log.WithPrefix(logger, severityKey, LevelDebug)
}

// Info returns a logger that includes a severity of level info
func Info(logger log.Logger) log.Logger {
	return log.WithPrefix(logger, severityKey, LevelInfo)
}

// Notice returns a logger that includes a severity of level notice
func Notice(logger log.Logger) log.Logger {
	return log.WithPrefix(logger, severityKey, LevelNotice)
}

// Warning returns a logger that includes a severity of level warning
func Warning(logger log.Logger) log.Logger {
	return log.WithPrefix(logger, severityKey, LevelWarning)
}

// Error returns a logger that includes a severity of level error
func Error(logger log.Logger) log.Logger {
	return log.WithPrefix(logger, severityKey, LevelError)
}

// Critical returns a logger that includes a severity of level critical
func Critical(logger log.Logger) log.Logger {
	return log.WithPrefix(logger, severityKey, LevelCritical)
}

// Alert returns a logger that includes a severity of level alert
func Alert(logger log.Logger) log.Logger {
	return log.WithPrefix(logger, severityKey, LevelAlert)
}

// Emergency returns a logger that includes a severity of level emergency
func Emergency(logger log.Logger) log.Logger {
	return log.WithPrefix(logger, severityKey, LevelEmergency)
}

type level byte

// LevelValue represents the value of the log level
type LevelValue struct {
	name string
	level
}

func (v *LevelValue) String() string { return v.name }
func (v *LevelValue) levelVal()      {}

const (
	levelDebug level = 1 << iota
	levelInfo
	levelNotice
	levelWarning
	levelError
	levelCritical
	levelAlert
	levelEmergency
)

var (
	LevelDebug     = &LevelValue{level: levelDebug, name: "debug"}
	LevelInfo      = &LevelValue{level: levelInfo, name: "info"}
	LevelNotice    = &LevelValue{level: levelNotice, name: "notice"}
	LevelWarning   = &LevelValue{level: levelWarning, name: "warning"}
	LevelError     = &LevelValue{level: levelError, name: "error"}
	LevelCritical  = &LevelValue{level: levelCritical, name: "critical"}
	LevelAlert     = &LevelValue{level: levelAlert, name: "alert"}
	LevelEmergency = &LevelValue{level: levelEmergency, name: "emergency"}
)

type filter struct {
	next     log.Logger
	minLevel *LevelValue
}

// LevelFromString returns the LevelValue represented by the given string.
// If the string doesn't represent a valid log level a default of LevelInfo is returned.
func LevelFromString(in string) *LevelValue {
	switch in {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "notice":
		return LevelNotice
	case "warning":
		return LevelWarning
	case "error":
		return LevelError
	case "critical":
		return LevelCritical
	case "alert":
		return LevelAlert
	case "emergency":
		return LevelEmergency
	}
	return LevelInfo
}

// NewFilter creates a new filtering loggering which surpresses all log entries
// with a severity below the specified LevelValue
func NewFilter(next log.Logger, minLevel *LevelValue) log.Logger {
	return &filter{
		next:     next,
		minLevel: minLevel,
	}
}

func (f *filter) Log(keyvals ...interface{}) error {
	for i := 0; i < len(keyvals); i += 2 {
		if stringVal, ok := keyvals[i].(string); ok {
			if stringVal == severityKey {
				if logLevel, ok := keyvals[i+1].(*LevelValue); ok {
					if logLevel.level >= f.minLevel.level {
						return f.next.Log(keyvals...)
					}
				}
			}
		}
	}
	return nil
}
