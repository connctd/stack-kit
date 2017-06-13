package logging

import (
	"github.com/go-kit/kit/log"
)

const (
	severityKey     = "severity"
	defaultErrorKey = "error"
)

var ErrorKey = defaultErrorKey

// SeverityAwareLogger tries to solve the problem that stackdriver logging regards
// every structured log entry as error if no severity is set. This logger wrapper inspects
// the key value pairs. If there is a key value pair with the key 'error' the severity will
// be set to error. By default severity will be info. If there is a key 'severity' this
// severity will be used.
type SeverityAwareLogger struct {
	log.Logger
}

func (s *SeverityAwareLogger) Log(keyvals ...interface{}) error {
	severity := LevelInfo
	severitySet := false
	errorIndex := -1
	for i := 0; i < len(keyvals); i += 2 {
		if stringVal, ok := keyvals[i].(string); ok {
			if stringVal == severityKey {
				severitySet = true
				break
			}

			if stringVal == ErrorKey {
				errorIndex = i
			}
		}
	}

	if errorIndex != -1 && !severitySet {
		if errorIndex+1 < len(keyvals) {
			switch v := keyvals[errorIndex+1].(type) {
			case string:
				if v != "nil" && v != "null" && v != "" {
					severity = LevelError
				} else {
					severity = LevelInfo
				}
			case error:
				severity = LevelError
			}
		} else {
			severity = LevelError
		}
	}

	if !severitySet {
		keyvals = append(keyvals, severityKey, severity)
	}
	return s.Logger.Log(keyvals...)
}
