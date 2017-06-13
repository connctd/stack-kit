package logging

import (
	"github.com/go-kit/kit/log"
)

type level byte

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
