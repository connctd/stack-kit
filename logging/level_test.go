package logging

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/go-kit/kit/log"

	"github.com/stretchr/testify/assert"
)

func TestFilterDebugLogging(t *testing.T) {
	logBuf := &bytes.Buffer{}
	logger := log.NewJSONLogger(logBuf)
	logger = &SeverityAwareLogger{logger}
	logger = NewFilter(logger, LevelInfo)

	logger.Log("msg", "something harmless", "value", 12, "bool", true, "severity", LevelDebug)
	assert.Empty(t, logBuf.String())
}

func TestFilterAllowDebugging(t *testing.T) {
	logBuf := &bytes.Buffer{}
	logger := log.NewJSONLogger(logBuf)
	logger = &SeverityAwareLogger{logger}
	logger = NewFilter(logger, LevelDebug)

	logger.Log("msg", "something harmless", "value", 12, "bool", true, "severity", LevelDebug)
	assert.Contains(t, logBuf.String(), `"severity":"debug"`)
}

type levelFunc func(log.Logger) log.Logger

func TestLevelHelpers(t *testing.T) {
	assert := assert.New(t)

	for levelVal, lf := range map[*LevelValue]levelFunc{
		LevelDebug:     Debug,
		LevelInfo:      Info,
		LevelNotice:    Notice,
		LevelWarning:   Warning,
		LevelError:     Error,
		LevelCritical:  Critical,
		LevelAlert:     Alert,
		LevelEmergency: Emergency,
	} {
		var logBuf bytes.Buffer
		logger := log.NewJSONLogger(&logBuf)
		lf(logger).Log("msg", "test")
		assert.Contains(logBuf.String(), fmt.Sprintf(`"severity":"%s"`, levelVal.String()))
	}
}

func TestLevelFromString(t *testing.T) {
	assert := assert.New(t)

	for stringVal, expectedVal := range map[string]*LevelValue{
		"debug":     LevelDebug,
		"info":      LevelInfo,
		"notice":    LevelNotice,
		"warning":   LevelWarning,
		"error":     LevelError,
		"critical":  LevelCritical,
		"alert":     LevelAlert,
		"emergency": LevelEmergency,
		"invalid":   LevelInfo,
	} {
		lv := LevelFromString(stringVal)
		assert.Equal(expectedVal, lv)
	}
}
