package logging

import (
	"bytes"
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
